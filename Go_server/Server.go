package main

import "fmt"

func Take_tokens()

type Server struct {
	Auth_key_word    string
	Refresh_key_word string
}

func main() {
	s := Server{"Auth_key_word", "Refresh_key_word"}
	fmt.Println(s)
}
