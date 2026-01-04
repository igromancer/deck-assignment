package queue

import (
	"context"
	"encoding/json"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ISender interface {
	Publish(ctx context.Context, msg data.JobPublic) error
}

type RabbitMQSender struct {
	Cfg             *config.Config
	QueueConnection *RabbitMQConnection
}

func (rs *RabbitMQSender) Publish(ctx context.Context, msg data.JobPublic) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return rs.QueueConnection.Ch.PublishWithContext(
		ctx,
		"",
		rs.Cfg.JobQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func NewSender(cfg *config.Config) (ISender, error) {
	connection := RabbitMQConnection{
		Cfg: cfg,
	}
	err := connection.Connect()
	s := RabbitMQSender{
		Cfg:             cfg,
		QueueConnection: &connection,
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
