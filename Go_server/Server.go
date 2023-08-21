package main

import (
	"context"
	"fmt"
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

func main() {
	s := Server{"Auth_key_word", "Refresh_key_word"}
	fmt.Println(s.Take_tokens(map[string]string{
		"user_id": "64dcd4c0aad456d0e90a9f3b",
	}))
	fmt.Println(s.Take_tokens(map[string]string{
		"refresh_token": "b'$2b$12$DiONzXtxVqWZmW.nTjMtIeHglhKPjCpB5CiR99XcCsvoZDtKrVA9S'",
	}))
}
