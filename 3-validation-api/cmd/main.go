package main

import (
	"net/http"
	"purple/links/configs"
	"purple/links/internal/verify"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, conf)

	server := http.Server{
		Addr:    conf.ServerConfig.Addr,
		Handler: router,
	}
	server.ListenAndServe()
}
