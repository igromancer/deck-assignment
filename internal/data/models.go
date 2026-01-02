package data

import (
	"gorm.io/gorm"
)

const (
	JobStatusPending    string = "pending"
	JobStatusProcessing string = "processing"
	JobStatusCompleted  string = "completed"
	JobStatusFailed     string = "failed"
)

type ApiKey struct {
	gorm.Model
	PublicId     string
	HashedSecret string
}

type Job struct {
	gorm.Model
	ApiKeyID       uint
	ApiKey         ApiKey
	Url            string
	Status         string
	ResultLocation string
}
