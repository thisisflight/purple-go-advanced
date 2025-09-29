package main

import "net/http"

func main() {
	router := http.NewServeMux()
	NewRandomHandler(router)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	server.ListenAndServe()
}
