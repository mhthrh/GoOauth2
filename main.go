package main

import (
	"GitHub.com/mhthrh/GoOauth/API"
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("run server")
	server := http.Server{
		Addr:    ":8585",
		Handler: &API.RequestHandler{},
	}
	server.ListenAndServe()

}
