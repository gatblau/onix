/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

const DbName = "doorman"

var (
	ErrDocumentAlreadyExists = errors.New("mongo: the document already exists")
	ErrDocumentNotFound      = errors.New("mongo: the document was not found")
)

// Db manage MongoDb connections
type Db struct {
	options *options.ClientOptions
}

func NewDb() *Db {
	return &Db{
		options: options.Client().ApplyURI(getDbConnString()),
	}
}

// ctx create a context with timeout of 30 seconds
func ctx() context.Context {
	context, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return context
}

// getDbConnString get the connection string to the MongoDb database
// e.g. mongodb://localhost:27017
// e.g. mongodb://user:password@127.0.0.1:27017/dbname?keepAlive=true&poolSize=30&autoReconnect=true&socketTimeoutMS=360000&connectTimeoutMS=360000
func getDbConnString() string {
	value := os.Getenv("DOORMAN_DB_CONN")
	if len(value) == 0 {
		panic("DOORMAN_DB_CONN not defined")
	}
	return value
}

// InsertObject insert a nameable object in the specified collection
func (db *Db) InsertObject(collection types.Collection, obj types.Nameable) (interface{}, error) {
	item, err := db.FindByName(collection, obj.GetName())
	if err != nil {
		return nil, err
	}
	// if the key was found
	if item.Err() != mongo.ErrNoDocuments {
		return nil, ErrDocumentAlreadyExists
	}
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(c)
	coll := client.Database(DbName).Collection(string(collection))
	// insert the key
	result, insertErr := coll.InsertOne(c, obj)
	if insertErr != nil {
		return nil, fmt.Errorf("cannot insert object into %s collection: %s", collection, err)
	}
	return result.InsertedID, nil
}

func (db *Db) UpdateObject(collection types.Collection, obj types.Nameable) (interface{}, error) {
	item, err := db.FindByName(collection, obj.GetName())
	if err != nil {
		return nil, err
	}
	// if the key was not found
	if item.Err() == mongo.ErrNoDocuments {
		return nil, ErrDocumentNotFound
	}
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(c)
	coll := client.Database(DbName).Collection(string(collection))
	result, updateErr := coll.UpdateOne(c, bson.M{"_id": obj.GetName()}, bson.M{"$set": obj})
	if updateErr != nil {
		return nil, updateErr
	}
	return result, nil
}

func (db *Db) UpsertObject(collection types.Collection, obj types.Nameable) (result interface{}, err error, resultCode int) {
	result, err = db.InsertObject(collection, obj)
	if err == nil {
		resultCode = 201
	} else if err == ErrDocumentAlreadyExists {
		// performs an update instead
		result, err = db.UpdateObject(collection, obj)
		if err != nil {
			resultCode = 500
			return nil, fmt.Errorf("cannot update document in collection %s with name %s: %s", collection, obj.GetName(), err), -1
		}
		resultCode = 200
	} else {
		resultCode = 500
	}
	return
}

// FindByName find an object by name
func (db *Db) FindByName(collection types.Collection, name string) (*mongo.SingleResult, error) {
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(c)
	coll := client.Database(DbName).Collection(string(collection))
	item := coll.FindOne(c, bson.M{"_id": name})
	return item, nil
}

// FindMany find a number of objects matching the specified filter
func (db *Db) FindMany(collection types.Collection, filter bson.M, query Query) error {
	if filter == nil {
		filter = bson.M{}
	}
	client, err := mongo.Connect(context.Background(), db.options)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	coll := client.Database(DbName).Collection(string(collection))
	cursor, findErr := coll.Find(context.Background(), filter)
	defer func() {
		// only closes the cursor if it exists
		if cursor != nil {
			cursor.Close(context.Background())
		}
	}()
	if findErr != nil {
		return findErr
	}
	return query(cursor)
}

type Query func(cursor *mongo.Cursor) error

func (db *Db) FindAll(collection types.Collection, results interface{}) error {
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return err
	}
	defer client.Disconnect(c)
	coll := client.Database(DbName).Collection(string(collection))
	cursor, err := coll.Find(c, make(map[string]interface{}))
	defer cursor.Close(c)
	if err != nil {
		return err
	}
	return cursor.All(c, &results)
}

// ObjectExists checks if an object exists in the specified collection
func (db *Db) ObjectExists(collection types.Collection, name string) bool {
	item, err := db.FindByName(collection, name)
	if err != nil {
		log.Printf("cannot retrieve item %s in collection %s: %s\n", name, collection, err)
		return false
	}
	return item.Err() != mongo.ErrNoDocuments
}
