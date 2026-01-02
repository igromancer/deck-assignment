package data

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Data access interface for jobs
type IJobRepository interface {
	Get(id uint) (*Job, error)
	Create(job *Job) (uint, error)
}

// Data access interface for api keys
type IApiKeyRepository interface {
	Get(ctx context.Context, publicId string) (*ApiKey, error)
	Create(ctx context.Context) (string, error)
}

// Implementation for accessing jobs
type JobRepository struct {
	Db *gorm.DB
}

func (j *JobRepository) Get(id uint) (*Job, error) {
	return &Job{}, nil
}

func (j *JobRepository) Create(job *Job) (uint, error) {
	return 1, nil
}

func NewJobrepository() (IJobRepository, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}
	jr := JobRepository{
		Db: db,
	}
	return &jr, nil
}

// Implementation for accessing api keys
type ApiKeyRepository struct {
	Db *gorm.DB
}

func (a *ApiKeyRepository) Create(ctx context.Context) (string, error) {
	// generate public and secret bits
	publicId := uuid.New().String()
	b := make([]byte, 32) // 256 bits
	rand.Read(b)
	secret := base64.RawURLEncoding.EncodeToString(b)
	apiKey := fmt.Sprintf("%s.%s", publicId, secret)
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// save publicId and hashedSecret in DB
	dbApiKey := ApiKey{
		PublicId:     publicId,
		HashedSecret: string(hashedSecret),
	}
	err = gorm.G[ApiKey](a.Db).Create(ctx, &dbApiKey)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

func (a *ApiKeyRepository) Get(ctx context.Context, publicId string) (*ApiKey, error) {
	apiKey, err := gorm.G[ApiKey](a.Db).Where("public_id = ?", publicId).First(ctx)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func NewApiKeyrepository() (IApiKeyRepository, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}
	r := ApiKeyRepository{
		Db: db,
	}
	return &r, nil
}
