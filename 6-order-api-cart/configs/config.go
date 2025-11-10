package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DBConfig
	SMTPConfig
	ServerConfig
	Auth AuthConfig
}

type DBConfig struct {
	Dsn string
}

type SMTPConfig struct {
	Email    string
	Password string
	Address  string
	Port     string
}

type ServerConfig struct {
	Schema string
	Addr   string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}
	return &Config{
		Db: DBConfig{
			Dsn: os.Getenv("DSN"),
		},
		SMTPConfig: SMTPConfig{
			Email:    os.Getenv("SMTP_EMAIL"),
			Password: os.Getenv("SMTP_PASSWORD"),
			Address:  os.Getenv("SMTP_ADDRESS"),
			Port:     os.Getenv("SMTP_PORT"),
		},
		ServerConfig: ServerConfig{
			Schema: os.Getenv("SCHEMA"),
			Addr:   os.Getenv("ADDR"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
