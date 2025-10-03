package main

import (
	"net/http"
	"purple/links/configs"
	"purple/links/files"
	"purple/links/internal/verify"
	"purple/links/pkg/db"
	"purple/links/storage"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDB(conf)
	db := files.NewJSONDB("storage.json")
	repo := storage.NewTokenRepository(db)
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, conf, repo)

	server := http.Server{
		Addr:    conf.ServerConfig.Addr,
		Handler: router,
	}
	server.ListenAndServe()
}
