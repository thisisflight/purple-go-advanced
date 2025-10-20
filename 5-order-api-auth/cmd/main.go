package main

import (
	"net/http"
	"purple/links/configs"
	"purple/links/internal/auth"
	"purple/links/internal/product"
	"purple/links/internal/session"
	"purple/links/internal/user"
	"purple/links/internal/verify"
	"purple/links/pkg/db"
	"purple/links/pkg/jwt"
	"purple/links/pkg/middleware"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetLevel(logrus.InfoLevel)
	conf := configs.LoadConfig()
	DB := db.NewDB(conf)
	jwtConfig := jwt.NewJWT(conf.Auth.Secret)
	// repositories
	productRepo := product.NewProductRepository(DB)
	userRepo := user.NewUserRepository(DB)
	sessionRepo := session.NewSessionRepository(DB)
	verifyRepo := verify.NewVerifyRepository(DB)
	productDeps := product.ProductHandlerDeps{ProductRepository: productRepo, Config: conf}
	authServiceDeps := auth.AuthServiceDeps{
		Conf:              conf,
		UserRepository:    userRepo,
		SessionRepository: sessionRepo,
		VerifyRepository:  verifyRepo,
		JWT:               jwtConfig}
	//services
	authService := auth.NewAuthService(authServiceDeps)
	authDeps := auth.AuthHandlerDeps{AuthService: authService}
	router := http.NewServeMux()
	// set handlers
	auth.NewAuthHandler(router, authDeps)
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
