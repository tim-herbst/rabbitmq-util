# rabbitmq-util
A common rabbitMQ go module

## Usage

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