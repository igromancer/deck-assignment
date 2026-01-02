package data

import (
	"fmt"

	"github.com/igromancer/deck-assignment/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
