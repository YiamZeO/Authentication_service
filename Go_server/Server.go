package main

import "fmt"

func Take_tokens(filter map[string]string) map[string]string {
	tokens := make(map[string]string)
	if _, ok := filter["_id"]; ok {

	} else if _, ok := filter["refresh"]; ok {

	} else {
		tokens["access_token"] = "Nil"
		tokens["refresh_token"] = "Nil"
	}
	return tokens
}

type Server struct {
	Auth_key_word    string
	Refresh_key_word string
}

func main() {
	s := Server{"Auth_key_word", "Refresh_key_word"}
	fmt.Println(s)
}
