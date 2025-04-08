package mongodb

import (
	"context"
	"fmt"

	"github.com/capigiba/capiary/internal/infra/db/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBAdapter[T any] struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewMongoDBAdapter[T any](client *mongo.Client, databaseName, collectionName string) *MongoDBAdapter[T] {
	ctx := context.Background()
	collection := client.Database(databaseName).Collection(collectionName)

	return &MongoDBAdapter[T]{
		collection: collection,
		ctx:        ctx,
	}
}

// ╔═════════════════════════════════════════╗
// ║     Updated FindWithQuery using Query   ║
// ╚═════════════════════════════════════════╝
func (m *MongoDBAdapter[T]) FindWithQuery(opts query.QueryOptions) ([]T, error) {
	// Build the filter and FindOptions via your new query package
	filter, findOpts := query.BuildMongoQuery(opts)
	// Then simply delegate to your existing Find method
	return m.Find(filter, findOpts)
}

// ╔═════════════════════════════════════════╗
// ║         Common MongoDB Functions        ║
// ╚═════════════════════════════════════════╝

func (m *MongoDBAdapter[T]) FindOne(filter interface{}) (*T, error) {
	var result T
	err := m.collection.FindOne(m.ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find one document: %v", err)
	}
	return &result, nil
}

func (m *MongoDBAdapter[T]) Find(filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := m.collection.Find(m.ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(m.ctx)

	var results []T
	for cursor.Next(m.ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

func (m *MongoDBAdapter[T]) UpdateOne(filter, update interface{}) error {
	_, err := m.collection.UpdateOne(m.ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}
	return nil
}

func (m *MongoDBAdapter[T]) BulkWrite(data map[string]T) error {
	var operations []mongo.WriteModel

	for key, value := range data {
		operations = append(operations,
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": key}).
				SetUpdate(bson.M{"$set": value}).
				SetUpsert(true))
	}

	_, err := m.collection.BulkWrite(m.ctx, operations)
	if err != nil {
		return fmt.Errorf("could not perform bulk write: %v", err)
	}

	return nil
}

func (m *MongoDBAdapter[T]) InsertOne(data T) (primitive.ObjectID, error) {
	result, err := m.collection.InsertOne(m.ctx, data)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to cast inserted ID to ObjectID")
	}

	return oid, nil
}
