package repository

import (
	"backend/internal/graphql/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBUserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBUserRepository(uri, dbName, collectionName string) (*MongoDBUserRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBUserRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *MongoDBUserRepository) UpsertUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Name == "" {
		// Insert new user
		result, err := r.collection.InsertOne(ctx, user)
		if err != nil {
			return models.User{}, err
		} else {
			user.ID = result.InsertedID.(primitive.ObjectID).Hex()
		}
	} else {
		// Update existing user
		filter := bson.M{"name": user.Name}
		update := bson.M{
			"$set": user,
		}
		opts := options.Update().SetUpsert(true)

		_, err := r.collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (r *MongoDBUserRepository) GetUsers() []*models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error getting users:", err)
		return nil
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			fmt.Println("Error decoding user:", err)
			continue
		}
		users = append(users, &user)
	}

	return users
}

func (r *MongoDBUserRepository) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoDBUserRepository) GetUserByName(name string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	fmt.Println("Name:", name)

	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
