package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct{
	Id primitive.ObjectID			`json:"id" bson:"_id,omitempty"`
	Name string						`json:"name" bson:"name"`	
	Amount string					`json:"amount" bson:"amount"`
	Status string					`json:"status" bson:"status"`
}