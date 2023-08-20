package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Err_check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	user_id := "64dcd4c0aad456d0e90a9f3b"
	base, err := url.Parse("http://127.0.0.1:5000/user/authentication")
	Err_check(err)
	base.RawQuery = url.Values{
		"user_id": {user_id},
	}.Encode()
	respons, err := http.PostForm(base.String(), nil)
	Err_check(err)
	defer respons.Body.Close()
	rec_data, err := io.ReadAll(respons.Body)
	Err_check(err)
	fmt.Println(string(rec_data))
}
