package jobs

import (
	"context"
	"net/url"
	"time"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
)

type IJobProcessor interface {
	ProcessJob(ctx context.Context, msg data.JobPublic) error
}

type IJobResultStore interface {
	WriteJobResult(jobID uint) error
}

type ScrapeJobProcessor struct {
	Cfg         *config.Config
	JobRepo     data.IJobRepository
	ResultStore IJobResultStore
	Sleep       func(time.Duration)
}

func (sjp *ScrapeJobProcessor) ProcessJob(ctx context.Context, msg data.JobPublic) error {
	sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusProcessing)
	sjp.Sleep(time.Second * 6)
	_, err := url.ParseRequestURI(msg.Url)
	if err != nil {
		sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusFailed)
		return err
	}
	err = sjp.ResultStore.WriteJobResult(msg.Id)
	if err != nil {
		sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusFailed)
		return err
	}
	sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusCompleted)
	return nil
}

func NewScrapeJobProcessor(cfg config.Config) (IJobProcessor, error) {
	repo, err := data.NewJobrepository(cfg)
	if err != nil {
		return nil, err
	}
	sjp := ScrapeJobProcessor{
		Cfg:         &cfg,
		JobRepo:     repo,
		ResultStore: NewJobResultStore(&cfg),
		Sleep:       time.Sleep,
	}
	return &sjp, nil
}
