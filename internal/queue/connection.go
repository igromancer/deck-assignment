package queue

import (
	"fmt"

	"github.com/igromancer/deck-assignment/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	Cfg  *config.Config
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func (rc *RabbitMQConnection) Connect() error {
	connUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rc.Cfg.RabbitMQUser,
		rc.Cfg.RabbitMQPassword,
		rc.Cfg.RabbitMQHost,
		rc.Cfg.RabbitMQPort,
	)
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return err
	}
	rc.Conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	rc.Ch = ch
	_, err = ch.QueueDeclare(
		rc.Cfg.JobQueueName,
		true,  // durable so messages survive broker restart
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		rc.Disconnect()
		return err
	}
	return nil
}

func (rc *RabbitMQConnection) Disconnect() {
	if rc.Ch != nil {
		_ = rc.Ch.Close()
		rc.Ch = nil
	}
	if rc.Conn != nil {
		_ = rc.Conn.Close()
		rc.Conn = nil
	}
}
