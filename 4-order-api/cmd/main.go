package main

import (
	"net/http"
	"purple/links/configs"
	"purple/links/files"
	"purple/links/internal/product"
	"purple/links/internal/verify"
	"purple/links/pkg/db"
	"purple/links/storage"
)

func main() {
	conf := configs.LoadConfig()
	fileDB := files.NewJSONDB("storage.json")
	mainDB := db.NewDB(conf)
	repo := storage.NewTokenRepository(fileDB)
	productRepo := product.NewProductRepository(mainDB)
	productDeps := product.ProductHandlerDeps{ProductRepository: productRepo}
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, conf, repo)
	product.NewProductHandler(router, productDeps)

	server := http.Server{
		Addr:    conf.ServerConfig.Addr,
		Handler: router,
	}
	server.ListenAndServe()
}
