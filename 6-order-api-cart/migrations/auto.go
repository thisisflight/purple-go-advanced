package main

import (
	"os"
	"purple/links/internal/order"
	"purple/links/internal/product"
	"purple/links/internal/session"
	"purple/links/internal/user"
	"purple/links/internal/verify"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&session.Session{})
	db.AutoMigrate(&verify.VerifyCode{})
	db.AutoMigrate(&order.Order{})
	db.AutoMigrate(&order.OrderItem{})
}
