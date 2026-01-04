package data

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func gormModel(id uint, createdAt, updatedAt time.Time) gorm.Model {
	return gorm.Model{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func TestToJobPublic(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 1, 10, 5, 0, 0, time.UTC)

	tests := []struct {
		name            string
		job             Job
		expectCompleted bool
	}{
		{
			name: "pending job",
			job: Job{
				Model:  gormModel(1, createdAt, updatedAt),
				Url:    "https://example.com",
				Status: JobStatusPending,
			},
			expectCompleted: false,
		},
		{
			name: "completed job",
			job: Job{
				Model:  gormModel(2, createdAt, updatedAt),
				Url:    "https://example.com",
				Status: JobStatusCompleted,
			},
			expectCompleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := ToJobPublic(&tt.job)

			assert.Equal(t, tt.job.ID, out.Id)
			assert.Equal(t, tt.job.Status, out.Status)
			assert.Equal(t, tt.job.Url, out.Url)
			assert.Equal(t, createdAt, out.CreatedAt)
			assert.Equal(t, "/jobs/"+strconv.Itoa(int(tt.job.ID)), out.StatusUrl)

			if tt.expectCompleted {
				assert.Equal(t, updatedAt, out.CompletedAt)
			} else {
				assert.True(t, out.CompletedAt.IsZero())
			}
		})
	}
}

func TestToJobResultPublic(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 1, 9, 10, 0, 0, time.UTC)

	t.Run("non-completed job does not set extracted_at", func(t *testing.T) {
		job := Job{
			Model:  gormModel(10, createdAt, updatedAt),
			Status: JobStatusProcessing,
		}

		out := ToJobResultPublic(&job, `{"foo":"bar"}`)

		assert.Equal(t, uint(10), out.Id)
		assert.Equal(t, JobStatusProcessing, out.Status)
		assert.Equal(t, createdAt, out.CreatedAt)
		assert.Equal(t, `{"foo":"bar"}`, out.Data)
		assert.True(t, out.ExtractedAt.IsZero())
	})

	t.Run("completed job sets extracted_at to UpdatedAt", func(t *testing.T) {
		job := Job{
			Model:  gormModel(11, createdAt, updatedAt),
			Status: JobStatusCompleted,
		}

		out := ToJobResultPublic(&job, `{"result":"ok"}`)

		assert.Equal(t, uint(11), out.Id)
		assert.Equal(t, JobStatusCompleted, out.Status)
		assert.Equal(t, createdAt, out.CreatedAt)
		assert.Equal(t, `{"result":"ok"}`, out.Data)
		assert.Equal(t, updatedAt, out.ExtractedAt)
	})
}
