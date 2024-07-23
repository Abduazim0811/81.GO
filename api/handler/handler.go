package handler

import (
	"net/http"

	"81.GO/internal/models"
	"81.GO/internal/mongodb"
	rabbitmq "81.GO/internal/rabbitmq/producer"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	db *mongodb.OrderMongodb
}

func NewOrderHandler(db *mongodb.OrderMongodb) *OrderHandler {
	return &OrderHandler{db: db}
}

func (o *OrderHandler) CreateOrders(c *gin.Context) {
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.Status = "Pending"

	err := rabbitmq.PublishOrder(order, "order.pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (o *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := o.db.StoreGetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) UpdateOrders(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order.Id = objectId

	err = rabbitmq.PublishOrder(order, "order.updated")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish updated order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order update request sent"})
}

func (o *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order := models.Order{Id: objectId}
	err = rabbitmq.PublishOrder(order, "order.deleted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order delete request sent"})
}
