package db

import (
	"context"
	"fmt"
	"github.com/egevorkyan/events/internal/event_store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ event_store.Storage = &db{}

type db struct {
	collection *mongo.Collection
}

func NewStorage(storage *mongo.Database, collection string) event_store.Storage {
	return &db{
		collection: storage.Collection(collection),
	}
}

func (s *db) Create(ctx context.Context, event event_store.EventLog) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.InsertOne(nCtx, event)
	if err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert objectid to hex")
}
