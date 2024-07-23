package consumer

import (
	"encoding/json"
	"log"

	"81.GO/internal/models"
	"81.GO/internal/mongodb"
	"github.com/streadway/amqp"
)

func StartConsumer(queueName string, handler func([]byte)) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func HandleCreateOrder(msg []byte) {
	var order models.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	mongodb, err := mongodb.NewOrder()
	if err != nil {
		log.Fatal(err)
	}

	if err = mongodb.StoreNewOrder(order); err != nil {
		log.Fatalf("Error saving task to MongoDB: %s", err)
	}
	log.Printf("Processed order: %s", order.Id)
}

func HandleUpdateOrder(msg []byte) {
	var order models.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	mongodb, err := mongodb.NewOrder()
	if err != nil {
		log.Fatal(err)
	}

	err = mongodb.StoreUpdateOrders(order.Id, order)
	if err != nil {
		log.Printf("Failed to update order: %s", err)
		return
	}

	log.Printf("Updated order: %s", order)
}

func HandleDeleteOrder(msg []byte) {
	var order models.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	mongodb, err := mongodb.NewOrder()
	if err != nil {
		log.Fatal(err)
	}

	err = mongodb.StoreDeleteOrders(order.Id)
	if err != nil {
		log.Printf("Failed to delete order: %s", err)
		return
	}

	log.Printf("Deleted order: %s", order.Id)
}
