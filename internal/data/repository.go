package data

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/igromancer/deck-assignment/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Data access interface for jobs
type IJobRepository interface {
	Get(ctx context.Context, id uint) (*Job, error)
	Create(ctx context.Context, job *Job) error
	List(ctx context.Context, apiKeyId uint, offset int, limit int) ([]Job, error)
	SetJobStatus(ctx context.Context, id uint, status string) error
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

func (j *JobRepository) Get(ctx context.Context, id uint) (*Job, error) {
	job, err := gorm.G[Job](j.Db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *JobRepository) Create(ctx context.Context, job *Job) error {
	err := gorm.G[Job](j.Db).Create(ctx, job)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobRepository) List(ctx context.Context, apiKeyId uint, offset int, limit int) ([]Job, error) {
	jobs, err := gorm.G[Job](j.Db).Where("api_key_id = ?", apiKeyId).Offset(offset).Limit(limit).Find(ctx)
	if err != nil {
		return []Job{}, err
	}
	return jobs, nil
}

func (j *JobRepository) SetJobStatus(ctx context.Context, id uint, status string) error {
	_, err := gorm.G[Job](j.Db).Where("id = ?", id).Update(ctx, "status", status)
	return err
}

func NewJobrepository(cfg config.Config) (IJobRepository, error) {
	db, err := GetConnection(cfg)
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

func NewApiKeyrepository(cfg config.Config) (IApiKeyRepository, error) {
	db, err := GetConnection(cfg)
	if err != nil {
		return nil, err
	}
	r := ApiKeyRepository{
		Db: db,
	}
	return &r, nil
}
