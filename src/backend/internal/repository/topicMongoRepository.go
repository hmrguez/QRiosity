package repository

import (
	"backend/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TopicMongoDBRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBTopicRepository(uri, dbName, collectionName string) (*TopicMongoDBRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &TopicMongoDBRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *TopicMongoDBRepository) GetAllTopics(ctx context.Context) ([]*domain.Topic, error) {
	var topics []*domain.Topic
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var topic domain.Topic
		if err := cursor.Decode(&topic); err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}
	return topics, nil
}

func (r *TopicMongoDBRepository) Insert(ctx context.Context, topics []*domain.Topic) error {
	var documents []interface{}
	for _, topic := range topics {
		documents = append(documents, topic)
	}
	_, err := r.collection.InsertMany(ctx, documents)
	return err
}
