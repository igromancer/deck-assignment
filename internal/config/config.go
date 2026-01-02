package config

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
}

func GetConfig() *Config {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load .env file: " + err.Error())
	}
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	return &c
}
