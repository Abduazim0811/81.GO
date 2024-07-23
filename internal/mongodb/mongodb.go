package mongodb

import (
	"context"
	"log"
	"time"

	"81.GO/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderMongodb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewOrder() (*OrderMongodb, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("Orders").Collection("order")
	return &OrderMongodb{client: client, collection: collection}, nil
}

func (o *OrderMongodb) StoreNewOrder(order models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.collection.InsertOne(ctx, order)
	return err
}

func (o *OrderMongodb) StoreGetOrders() ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []*models.Order
	cursor, err := o.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(ctx, &orders); err != nil {
		log.Fatal(err)
	}
	return orders, nil
}

func (o *OrderMongodb) StoreUpdateOrders(id primitive.ObjectID, order models.Order) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.collection.UpdateByID(ctx, id, bson.M{"$set": order})
	return err
}

func (o *OrderMongodb) StoreDeleteOrders(id primitive.ObjectID) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
