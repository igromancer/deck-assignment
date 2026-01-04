package jobs

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
)

type IJobProcessor interface {
	ProcessJob(ctx context.Context, msg data.JobPublic) error
}

type ScrapeJobProcessor struct {
	JobRepo data.IJobRepository
}

func (sjp *ScrapeJobProcessor) ProcessJob(ctx context.Context, msg data.JobPublic) error {
	fmt.Println("================ Processing a job ===================")
	fmt.Println(msg)

	sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusProcessing)
	time.Sleep(time.Second * 6)
	_, err := url.ParseRequestURI(msg.Url)
	if err != nil {
		sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusFailed)
		return err
	}
	// Write dummy string as job result
	sjp.JobRepo.SetJobStatus(ctx, msg.Id, data.JobStatusCompleted)
	return nil
}

func NewScrapeJobProcessor(cfg config.Config) (IJobProcessor, error) {
	repo, err := data.NewJobrepository(cfg)
	if err != nil {
		return nil, err
	}
	sjp := ScrapeJobProcessor{
		JobRepo: repo,
	}
	return &sjp, nil
}
