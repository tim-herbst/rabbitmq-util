package consumer

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	conn *amqp.Connection
}

func NewConsumer(rabbitMQURL string) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return &Consumer{conn: conn}, nil
}

func (c *Consumer) Consume(queue string, handler func(delivery amqp.Delivery)) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Printf("failed to close channel: %v", err)
		}
	}(ch)
	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	for d := range deliveries {
		log.Println("Received a message: ", string(d.Body))
		handler(d)
	}
	return nil
}
