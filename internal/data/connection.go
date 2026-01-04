package data

import (
	"fmt"

	"github.com/igromancer/deck-assignment/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
}

func GetConnection(cfg config.Config) (*gorm.DB, error) {
	dsn := PostgresDSN(cfg)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
