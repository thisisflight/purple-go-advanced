package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPConfig
	ServerConfig
}

type SMTPConfig struct {
	Email    string
	Password string
	Address  string
}

type ServerConfig struct {
	Addr string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}
	return &Config{
		SMTPConfig: SMTPConfig{
			Email:    os.Getenv("SMTP_EMAIL"),
			Password: os.Getenv("SMTP_PASSWORD"),
			Address:  os.Getenv("SMTP_ADDRESS"),
		},
		ServerConfig: ServerConfig{
			Addr: os.Getenv("ADDR"),
		},
	}
}
