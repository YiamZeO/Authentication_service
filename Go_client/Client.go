package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func Err_check(err error) {
	if err != nil {
		panic(err)
	}
}

type User struct {
	User_id string
	Tokens  map[string]string
}

func (u *User) Auth(a_url string) {
	base, err := url.Parse(a_url)
	Err_check(err)
	base.RawQuery = url.Values{
		"user_id": {u.User_id},
	}.Encode()
	response, err := http.Post(base.String(), "application/x-www-form-urlencoded", nil)
	Err_check(err)
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&u.Tokens)
	Err_check(err)
}

func (u *User) Refresh(a_url string) {
	data := make([]byte, base64.StdEncoding.EncodedLen(len(u.Tokens["refresh_token"])))
	base64.StdEncoding.Encode(data, []byte(u.Tokens["refresh_token"]))
	response, err := http.Post(a_url, "application/x-www-form-urlencoded", bytes.NewBuffer(data))
	Err_check(err)
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&u.Tokens)
	Err_check(err)
}

func main() {

	//Rest маршруты
	urls := map[string]string{
		"Auth":    "http://127.0.0.1:5000/user/authentication",
		"Refresh": "http://127.0.0.1:5000/user/refresh",
	}

	fmt.Println("Enter user id: ")
	var user_id string
	fmt.Scanln(&user_id)

	// Пользователь: его id, токены
	u := User{
		User_id: user_id,
		Tokens: map[string]string{
			"access_token":  "Nil",
			"refresh_token": "Nil",
		},
	}
	for loop := true; loop; {
		fmt.Print("\033[H\033[2J")
		fmt.Printf("User id: %s\n", u.User_id)
		fmt.Printf("Access token: %s\n", u.Tokens["access_token"])
		fmt.Printf("Refresh token: %s\n\n", u.Tokens["refresh_token"])
		fmt.Printf("1. Authenticate\n")
		fmt.Printf("2. Refresh\n")
		fmt.Printf("3. Exit\n")
		fmt.Printf("Enter your choice: ")
		var choice int
		fmt.Scanf("%d", &choice)
		switch choice {
		case 1:
			u.Auth(urls["Auth"])
		case 2:
			u.Refresh(urls["Refresh"])
		case 3:
			loop = false
		default:
			fmt.Printf("Invalid choice\n")
		}
	}
}
