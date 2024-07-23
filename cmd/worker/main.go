package main

import (
	"81.GO/internal/rabbitmq/consumer"
)

func main() {
	go consumer.StartConsumer("order.pending", consumer.HandleCreateOrder)
	go consumer.StartConsumer("order.updated", consumer.HandleUpdateOrder)
	go consumer.StartConsumer("order.deleted", consumer.HandleDeleteOrder)

	select {}
}
