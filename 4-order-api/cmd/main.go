package main

import (
	"net/http"
	"purple/links/configs"
	"purple/links/files"
	"purple/links/internal/product"
	"purple/links/internal/verify"
	"purple/links/pkg/db"
	"purple/links/pkg/middleware"
	"purple/links/storage"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetLevel(logrus.InfoLevel)
	conf := configs.LoadConfig()
	fileDB := files.NewJSONDB("storage.json")
	mainDB := db.NewDB(conf)
	repo := storage.NewTokenRepository(fileDB)
	productRepo := product.NewProductRepository(mainDB)
	productDeps := product.ProductHandlerDeps{ProductRepository: productRepo}
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, conf, repo)
	product.NewProductHandler(router, productDeps)

	middlewareChain := middleware.Chain(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    conf.ServerConfig.Addr,
		Handler: middlewareChain(router),
	}
	server.ListenAndServe()
}
