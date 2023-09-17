package publisher

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn *amqp.Connection
}

func NewPublisher(rabbitMQURL string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Publish(exchange, routingKey string, msg []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Printf("failed to close channel: %v", err)
		}
	}(ch)

	err = ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	return nil
}
