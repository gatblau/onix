package core

/*
  Onix Config Manager - MongoDb event receiver for Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"context"
	"github.com/gatblau/onix/pilotctl/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Db manage MongoDb connections
type Db struct {
	options *options.ClientOptions
}

func NewDb() (*Db, error) {
	return &Db{
		options: options.Client().ApplyURI(getDbConnString()),
	}, nil
}

// Events return the events collection
func (db *Db) Events(client *mongo.Client) *mongo.Collection {
	return client.Database("syslog").Collection("events")
}

// Insert events in the data store
func (db *Db) Insert(events *types.Events) error {
	// convert input events to []interface{}
	var documents []interface{}
	for _, ev := range events.Events {
		documents = append(documents, ev)
	}
	client, err := mongo.Connect(ctx(), db.options)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx())

	// insert payload in events collection
	_, err = db.Events(client).InsertMany(ctx(), documents)
	return err
}

func (db *Db) Query(filter bson.M) (*types.Events, error) {
	client, err := mongo.Connect(ctx(), db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx())
	// create a cursor over the whole collection
	cursor, err := db.Events(client).Find(ctx(), filter)
	if err != nil {
		return nil, err
	}
	var results []types.Event
	// return elements that match the criteria
	if err = cursor.All(ctx(), &results); err != nil {
		return nil, err
	}
	return &types.Events{Events: results}, nil
}

// ctx create a context with timeout of 10 seconds
func ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
