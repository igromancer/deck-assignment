package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ISender interface {
	Publish(ctx context.Context, msg data.JobPublic) error
	Connect() error
	Disconnect()
}

type RabbitMQSender struct {
	Cfg  *config.Config
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func (rs *RabbitMQSender) Connect() error {
	connUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rs.Cfg.RabbitMQUser,
		rs.Cfg.RabbitMQPassword,
		rs.Cfg.RabbitMQHost,
		rs.Cfg.RabbitMQPort,
	)
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return err
	}
	rs.Conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	rs.Ch = ch
	_, err = ch.QueueDeclare(
		rs.Cfg.JobQueueName,
		true,  // durable so messages survive broker restart
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		rs.Disconnect()
		return err
	}
	return nil
}

func (rs *RabbitMQSender) Disconnect() {
	if rs.Ch != nil {
		_ = rs.Ch.Close()
		rs.Ch = nil
	}
	if rs.Conn != nil {
		_ = rs.Conn.Close()
		rs.Conn = nil
	}
}

func (rs *RabbitMQSender) Publish(ctx context.Context, msg data.JobPublic) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return rs.Ch.PublishWithContext(
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
	s := RabbitMQSender{
		Cfg: cfg,
	}
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	return &s, nil
}
