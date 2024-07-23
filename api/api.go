package api

import (
	"log"

	"81.GO/api/handler"
	"81.GO/internal/mongodb"
	rabbitmq "81.GO/internal/rabbitmq/producer"
	"github.com/gin-gonic/gin"
)

func Routes() {
	router := gin.Default()

	db, err := mongodb.NewOrder()
	if err != nil {
		log.Fatal(err)
	}

	rabbitmq.InitRabbitMQ()

	orderhandler := handler.NewOrderHandler(db)
	router.POST("/orders", orderhandler.CreateOrders)
	router.GET("/orders/:id", orderhandler.GetOrders)
	router.PUT("/orders/:id", orderhandler.UpdateOrders)
	router.DELETE("/orders/:id", orderhandler.DeleteOrder)

	router.Run(":8888")
}
