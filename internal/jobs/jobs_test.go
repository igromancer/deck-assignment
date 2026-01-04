package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/igromancer/deck-assignment/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJobRepo struct{ mock.Mock }

func (m *MockJobRepo) Get(ctx context.Context, id uint) (*data.Job, error) {
	args := m.Called(ctx, id)
	job, _ := args.Get(0).(*data.Job)
	return job, args.Error(1)
}

func (m *MockJobRepo) Create(ctx context.Context, job *data.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobRepo) List(ctx context.Context, apiKeyId uint, offset int, limit int) ([]data.Job, error) {
	args := m.Called(ctx, apiKeyId, offset, limit)
	jobs, _ := args.Get(0).([]data.Job)
	return jobs, args.Error(1)
}

func (m *MockJobRepo) SetJobStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockResultStore struct{ mock.Mock }

func (m *MockResultStore) WriteJobResult(jobID uint) error {
	args := m.Called(jobID)
	return args.Error(0)
}

func TestProcessJob_Success(t *testing.T) {
	repo := new(MockJobRepo)
	store := new(MockResultStore)

	p := &ScrapeJobProcessor{
		JobRepo:     repo,
		ResultStore: store,
		Sleep:       func(time.Duration) {}, // no-op
	}

	ctx := context.Background()
	job := data.JobPublic{
		Id:  1,
		Url: "https://example.com",
	}

	repo.On("SetJobStatus", ctx, uint(1), data.JobStatusProcessing).Return(nil).Once()
	store.On("WriteJobResult", uint(1)).Return(nil).Once()
	repo.On("SetJobStatus", ctx, uint(1), data.JobStatusCompleted).Return(nil).Once()

	err := p.ProcessJob(ctx, job)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	store.AssertExpectations(t)
}
