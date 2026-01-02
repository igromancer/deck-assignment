package data

import (
	"gorm.io/gorm"
)

type ApiKey struct {
	gorm.Model
	PublicId     string
	HashedSecret string
}

type Job struct {
	gorm.Model
	ApiKeyID       int
	ApiKey         ApiKey
	Url            string
	Status         string
	ResultLocation string
}
