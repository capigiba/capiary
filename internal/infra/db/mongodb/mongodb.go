package mongodb

import (
	"context"
	"sync"

	"github.com/capigiba/capiary/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client *mongo.Client
}

// Singleton instance and mutex
var (
	mongoOnce       sync.Once
	mongoDBInstance *MongoDBClient
)

var log = logger.NewLogger("mongodb")

// NewMongoDBClient creates a new MongoDBClient instance if not already created (Singleton)
func NewMongoDBClient(mongoDBEndPoint string) *MongoDBClient {
	mongoOnce.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBEndPoint))
		if err != nil {
			log.Errorf("Failed to connect to MongoDB: %v", err)
		}

		mongoDBInstance = &MongoDBClient{
			client: client,
		}
	})

	return mongoDBInstance
}

// GetClient returns the underlying *mongo.Client
func (m *MongoDBClient) GetClient() *mongo.Client {
	return m.client
}

// Close closes the MongoDB client connection when needed
func (m *MongoDBClient) Close() error {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		log.Errorf("Error closing MongoDB client: %v", err)
		return err
	}
	log.Infof("MongoDB client closed successfully.")
	return nil
}
