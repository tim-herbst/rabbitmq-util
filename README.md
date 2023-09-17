# rabbitmq-util

A common rabbitMQ go module

## Binder

The binder is a tool to bind exchanges and queues. It can be used to set up the rabbitMQ environment for a microservice.
Create a configuration file which holds the binding information as YAML format. For example:

```yaml
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
bindings:
  - exchange: "user"
    queue: "user"
    routingKey: "user"
    single_active_consumer: true
```

Then, create a rabbitMQ client with the configuration file:

```go
package main

import (
	"log"
	"rabbitmq-util/pkg/binder"
)

func main() {
	config, err := binder.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = binder.SetupBindings(config.RabbitMQ.URL, config.Bindings)
	if err != nil {
		log.Fatalf("Failed to setup bindings: %v", err)
	}

	log.Println("Bindings are set up successfully")
}

```

## Publisher and Consumer

The publisher and Consumer is a tool to publish and consume messages to and from rabbitMQ. It can be used to send and receive messages to rabbitMQ in a microservice.

```go
package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rabbitmq-util/pkg/binder"
	"rabbitmq-util/pkg/consumer"
	"rabbitmq-util/pkg/publisher"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func main() {
	config, err := binder.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	c, _ := consumer.NewConsumer(config.RabbitMQ.URL)
	err = c.Consume("queue", func(d amqp.Delivery) {
		fmt.Println(string(d.Body))
	})
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	m := Message {
		ID:      "1",
        Content: "Hello World",
    }
	p, _ := publisher.NewPublisher(config.RabbitMQ.URL)
	err = p.Publish("exchange", "routingKey", m)
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
}

```