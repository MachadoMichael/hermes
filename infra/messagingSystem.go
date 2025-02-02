package infra

import (
	"log"

	"github.com/streadway/amqp"
)

type MSConfig struct {
	URL string
}

type MSClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewMSClient(config MSConfig) (*MSClient, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &MSClient{conn: conn, ch: ch}, nil
}

func (c *MSClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *MSClient) Publish(queueName string, body []byte) error {
	_, err := c.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = c.ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Message published to queue: %s", queueName)
	return nil
}

func (c *MSClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	_, err := c.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, err
	}

	msgs, err := c.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
