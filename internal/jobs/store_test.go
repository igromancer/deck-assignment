package jobs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockFS struct{ mock.Mock }

func (m *MockFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockFS) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	b, _ := args.Get(0).([]byte)
	return b, args.Error(1)
}

func TestJobResultStore_getJobFileName(t *testing.T) {
	s := &JobResultStore{
		Cfg: &config.Config{JobResultPath: "/tmp"},
		Fs:  new(MockFS),
	}

	assert.Equal(t, "scrape-job-1.json", s.getJobFileName(1))
	assert.Equal(t, "scrape-job-42.json", s.getJobFileName(42))
}

func TestJobResultStore_WriteJobResult_Success(t *testing.T) {
	fs := new(MockFS)
	cfg := &config.Config{JobResultPath: "/data/results"}
	s := &JobResultStore{Cfg: cfg, Fs: fs}

	jobID := uint(7)
	expectedPath := filepath.Join(cfg.JobResultPath, "scrape-job-7.json")
	expectedContent := []byte(`{ "title": "Dummy", "content": "Dummy" }`)

	fs.On("WriteFile", expectedPath, expectedContent, os.FileMode(0644)).Return(nil).Once()

	err := s.WriteJobResult(jobID)
	require.NoError(t, err)
	fs.AssertExpectations(t)
}

func TestJobResultStore_WriteJobResult_Error(t *testing.T) {
	fs := new(MockFS)
	cfg := &config.Config{JobResultPath: "/data/results"}
	s := &JobResultStore{Cfg: cfg, Fs: fs}

	jobID := uint(7)
	expectedPath := filepath.Join(cfg.JobResultPath, "scrape-job-7.json")

	fs.On("WriteFile", expectedPath, mock.Anything, os.FileMode(0644)).
		Return(assert.AnError).
		Once()

	err := s.WriteJobResult(jobID)
	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
	fs.AssertExpectations(t)
}

func TestJobResultStore_ReadJobResult_Success(t *testing.T) {
	fs := new(MockFS)
	cfg := &config.Config{JobResultPath: "/data/results"}
	s := &JobResultStore{Cfg: cfg, Fs: fs}

	jobID := uint(99)
	expectedPath := filepath.Join(cfg.JobResultPath, "scrape-job-99.json")
	expectedBytes := []byte(`{"ok":true}`)

	fs.On("ReadFile", expectedPath).Return(expectedBytes, nil).Once()

	out, err := s.ReadJobResult(jobID)
	require.NoError(t, err)
	assert.Equal(t, `{"ok":true}`, out)
	fs.AssertExpectations(t)
}

func TestJobResultStore_ReadJobResult_Error(t *testing.T) {
	fs := new(MockFS)
	cfg := &config.Config{JobResultPath: "/data/results"}
	s := &JobResultStore{Cfg: cfg, Fs: fs}

	jobID := uint(99)
	expectedPath := filepath.Join(cfg.JobResultPath, "scrape-job-99.json")

	fs.On("ReadFile", expectedPath).Return([]byte(nil), assert.AnError).Once()

	out, err := s.ReadJobResult(jobID)
	require.Error(t, err)
	assert.Equal(t, "", out)
	assert.ErrorIs(t, err, assert.AnError)
	fs.AssertExpectations(t)
}
