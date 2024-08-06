package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoUser struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}
