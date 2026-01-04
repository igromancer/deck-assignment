package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig_ReadsEnv(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_NAME", "deck")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")

	t.Setenv("RABBITMQ_USER", "guest")
	t.Setenv("RABBITMQ_PASSWORD", "guest")
	t.Setenv("RABBITMQ_HOST", "rabbitmq")
	t.Setenv("RABBITMQ_PORT", "5672")

	t.Setenv("JOB_QUEUE_NAME", "jobs")
	t.Setenv("JOB_RESULT_PATH", "/data/results")

	cfg := GetConfig()
	require.NotNil(t, cfg)

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "deck", cfg.DBName)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)

	assert.Equal(t, "guest", cfg.RabbitMQUser)
	assert.Equal(t, "guest", cfg.RabbitMQPassword)
	assert.Equal(t, "rabbitmq", cfg.RabbitMQHost)
	assert.Equal(t, "5672", cfg.RabbitMQPort)

	assert.Equal(t, "jobs", cfg.JobQueueName)
	assert.Equal(t, "/data/results", cfg.JobResultPath)

	// ensure no unexpected env mutation happened (sanity)
	_, ok := os.LookupEnv("DB_HOST")
	assert.True(t, ok)
}
