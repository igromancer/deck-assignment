package data

import (
	"testing"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDSN(t *testing.T) {
	cfg := config.Config{
		DBHost:     "localhost",
		DBUser:     "postgres",
		DBPassword: "secret",
		DBName:     "deck",
		DBPort:     5432,
	}

	dsn := PostgresDSN(cfg)

	// exactly matches your fmt string
	assert.Equal(t, "host=localhost user=postgres password=secret dbname=deck port=5432", dsn)
}
