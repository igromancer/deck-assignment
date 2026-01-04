package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
	"github.com/igromancer/deck-assignment/internal/jobs"
)

type IJobReceiver interface {
	Run(ctx context.Context) error
}

type JobReceiver struct {
	Cfg             *config.Config
	QueueConnection *RabbitMQConnection
	Processor       jobs.IJobProcessor
}

func (jr *JobReceiver) Run(ctx context.Context) error {
	msgs, err := jr.QueueConnection.Ch.Consume(
		jr.Cfg.JobQueueName,
		"job-receiver-1", // consumer tag
		false,            // turn off auto ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return errors.New("delivery channel closed")
			}
			var jobMessage data.JobPublic
			if err := json.Unmarshal(d.Body, &jobMessage); err != nil {
				// reject without requeue (poison message)
				_ = d.Reject(false)
				continue
			}
			if err = jr.Processor.ProcessJob(ctx, jobMessage); err != nil {
				_ = d.Nack(false, true)
				continue
			}
			_ = d.Ack(false)
		}
	}
}

func NewJobReceiver() (IJobReceiver, error) {
	cfg := config.GetConfig()
	connection := RabbitMQConnection{
		Cfg: cfg,
	}
	err := connection.Connect()
	if err != nil {
		return nil, err
	}
	sjp, err := jobs.NewScrapeJobProcessor(*cfg)
	if err != nil {
		return nil, err
	}
	jr := JobReceiver{
		Cfg:             cfg,
		QueueConnection: &connection,
		Processor:       sjp,
	}
	return &jr, nil
}
