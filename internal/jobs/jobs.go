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

type ScrapeJobProcessor struct {
	Cfg     *config.Config
	JobRepo data.IJobRepository
}

func (sjp *ScrapeJobProcessor) ProcessJob(ctx context.Context, msg data.JobPublic) error {
	sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusProcessing)
	time.Sleep(time.Second * 6)
	_, err := url.ParseRequestURI(msg.Url)
	if err != nil {
		sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusFailed)
		return err
	}
	store := NewJobResultStore(sjp.Cfg)
	err = store.WriteJobResult(msg.Id)
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
		Cfg:     &cfg,
		JobRepo: repo,
	}
	return &sjp, nil
}
