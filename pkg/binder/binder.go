package binder

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

type Config struct {
	RabbitMQ RabbitMQConfig  `yaml:"rabbitmq"`
	Bindings []BindingConfig `yaml:"bindings"`
}

type RabbitMQConfig struct {
	URL string `yaml:"url"`
}

type BindingConfig struct {
	Exchange             string `yaml:"exchange"`
	Queue                string `yaml:"queue"`
	RoutingKey           string `yaml:"routing_key"`
	SingleActiveConsumer bool   `yaml:"single_active_consumer"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SetupBindings(rabbitMQURL string, bindings []BindingConfig) error {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %s", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Printf("failed to close channel: %s", err)
		}
	}(ch)

	for _, binding := range bindings {
		err = ch.ExchangeDeclarePassive(binding.Exchange, "topic", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("exchange %s does not exists: %w", binding.Exchange, err)
		}

		_, err = ch.QueueDeclare(binding.Queue, true, false, false, false, amqp.Table{
			"x-single-active-consumer": binding.SingleActiveConsumer,
		})
		if err != nil {
			return fmt.Errorf("failed to declare a queue: %w", err)
		}

		err = ch.QueueBind(binding.Queue, binding.RoutingKey, binding.Exchange, false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind a queue: %w", err)
		}
		log.Printf("bound queue: %s to exchange: %s with routing key: %s", binding.Queue, binding.Exchange, binding.RoutingKey)
	}
	return nil
}
