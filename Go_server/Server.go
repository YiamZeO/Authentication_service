package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Err_check(err error) {
	if err != nil {
		panic(err)
	}
}

type Claims struct {
	User_id string `json:"user_id"`
	jwt.RegisteredClaims
}

type Server struct {
	Auth_key_word    string
	Refresh_key_word string
}

func (s *Server) Take_tokens(filter map[string]string) map[string]string {
	tokens := map[string]string{
		"access_token":  "Nil",
		"refresh_token": "Nil",
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	Err_check(err)
	defer client.Disconnect(context.TODO())
	users_coll := client.Database("Authentication_bd").Collection("Users")
	m_filter := bson.M{}
	if _, ok := filter["user_id"]; ok {
		m_filter["_id"], err = primitive.ObjectIDFromHex(filter["user_id"])
		Err_check(err)
	} else if _, ok := filter["refresh_token"]; ok {
		m_filter["refresh_token"] = filter["refresh_token"]
	} else {
		return tokens
	}
	c_user := map[string]string{}
	err = users_coll.FindOne(context.TODO(), m_filter).Decode(&c_user)
	Err_check(err)
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		User_id: c_user["_id"],
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	auth_token, err := token.SignedString([]byte(s.Auth_key_word))
	Err_check(err)
	expirationTime = time.Now().Add(180 * 24 * time.Hour)
	refresh_token_str := s.Refresh_key_word + expirationTime.Format(time.DateTime)
	refresh_token_bcrypt, err := bcrypt.GenerateFromPassword([]byte(refresh_token_str), bcrypt.DefaultCost)
	Err_check(err)
	tokens["access_token"] = auth_token
	tokens["refresh_token"] = string(refresh_token_bcrypt)
	mongo_user_up := bson.M{
		"$set": bson.M{
			"refresh_token": tokens["refresh_token"],
		},
	}
	_, err = users_coll.UpdateOne(context.TODO(), m_filter, mongo_user_up)
	Err_check(err)
	return tokens
}

func (s *Server) Post_Auth(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("user_id")
	filter := map[string]string{"user_id": user_id}
	tokens := s.Take_tokens(filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
	fmt.Println("User " + r.RemoteAddr + " get post auth operation")
}

func (s *Server) Post_Refresh(w http.ResponseWriter, r *http.Request) {
	body_data, err := io.ReadAll(r.Body)
	Err_check(err)
	data := make([]byte, base64.StdEncoding.DecodedLen(len(body_data)))
	_, err = base64.StdEncoding.Decode(data, body_data)
	Err_check(err)
	data_str := string(data)
	filter := map[string]string{"refresh_token": data_str}
	tokens := s.Take_tokens(filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
	fmt.Println("User " + r.RemoteAddr + " get post refresh operation")
}

func main() {

	// Структура для сервера
	server := Server{
		Auth_key_word:    "Auth_key_word",    // Secret word для токена авторизации
		Refresh_key_word: "Refresh_key_word", // Secret word для refresh токена
	}
	http.HandleFunc("/user/authentication", server.Post_Auth)
	http.HandleFunc("/user/refresh", server.Post_Refresh)
	fmt.Println("Server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
