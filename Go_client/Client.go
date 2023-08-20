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
	base, err := url.Parse("http://127.0.0.1:5000/user/authentication")
	Err_check(err)
	params := url.Values{}
	params.Add("user_id", "64dcd4c0aad456d0e90a9f3b")
	base.RawQuery = params.Encode()
	respons, err := http.Get(base.String())
	Err_check(err)
	defer respons.Body.Close()
	rec_data, err := io.ReadAll(respons.Body)
	Err_check(err)
	fmt.Println(string(rec_data))
}
