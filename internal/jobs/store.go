package jobs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/igromancer/deck-assignment/internal/config"
)

type JobResultStore struct {
	Cfg *config.Config
	Fs  IFs
}

func (s *JobResultStore) getJobFileName(jobID uint) string {
	return fmt.Sprintf("scrape-job-%v.json", jobID)
}

func (s *JobResultStore) WriteJobResult(jobID uint) error {
	result := `{ "title": "Dummy", "content": "Dummy" }`
	fileName := filepath.Join(s.Cfg.JobResultPath, s.getJobFileName(jobID))
	if err := s.Fs.WriteFile(fileName, []byte(result), 0644); err != nil {
		return err
	}
	return nil
}

func (s *JobResultStore) ReadJobResult(jobID uint) (string, error) {
	fileName := filepath.Join(s.Cfg.JobResultPath, s.getJobFileName(jobID))
	data, err := s.Fs.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewJobResultStore(cfg *config.Config) *JobResultStore {
	store := JobResultStore{
		Cfg: cfg,
		Fs:  &FuncFS{},
	}
	return &store
}

type IFs interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
}

type FuncFS struct{}

func (f *FuncFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (f *FuncFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
