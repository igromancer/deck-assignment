package data

import (
	"fmt"
	"time"

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

type JobPublic struct {
	Id          uint      `json:"job_id"`
	Status      string    `json:"status"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	StatusUrl   string    `json:"status_url"`
}

func ToJobPublic(j *Job) JobPublic {
	pj := JobPublic{
		Id:        j.ID,
		Status:    j.Status,
		Url:       j.Url,
		CreatedAt: j.CreatedAt,
		StatusUrl: fmt.Sprintf("/jobs/%v", j.ID),
	}
	if j.Status == JobStatusCompleted {
		pj.CompletedAt = j.UpdatedAt
	}
	return pj
}

type JobResultPublic struct {
	Id          uint           `json:"job_id"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ExtractedAt time.Time      `json:"extracted_at"`
	Data        map[string]any `json:"data"`
}

func ToJobResultPublic(j *Job, data map[string]any) JobResultPublic {
	jr := JobResultPublic{
		Id:        j.ID,
		Status:    j.Status,
		CreatedAt: j.CreatedAt,
		Data:      data,
	}
	if j.Status == JobStatusCompleted {
		jr.ExtractedAt = j.UpdatedAt
	}
	return jr
}
