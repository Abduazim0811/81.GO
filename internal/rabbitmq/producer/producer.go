package rabbitmq

import (
	"encoding/json"
	"log"

	"81.GO/internal/models"
	"github.com/streadway/amqp"
)

var channel *amqp.Channel

func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

    ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

    err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
	}
}

func PublishOrder(order models.Order, routingKey string) error {
	
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"orders_exchange",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
