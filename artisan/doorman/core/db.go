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
		return nil, fmt.Errorf("object in %s collection with name %s already exist", collection, obj.GetName())
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
func (db *Db) FindMany(collection types.Collection, filter bson.M, results interface{}) error {
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return err
	}
	defer client.Disconnect(c)
	coll := client.Database(DbName).Collection(string(collection))
	cursor, findErr := coll.Find(ctx(), filter)
	if findErr != nil {
		return findErr
	}
	// return elements that match the criteria
	if err = cursor.All(ctx(), &results); err != nil {
		return err
	}
	return nil
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

// FindKeys retrieves one or more keys matching the specifies criteria decrypting the value of any private key
func (db *Db) FindKeys(filter bson.M) ([]types.Key, error) {
	var results []types.Key
	err := db.FindMany(types.KeysColl, filter, results)
	if err != nil {
		return nil, err
	}
	for i, key := range results {
		if key.IsPrivate {
			dec, decErr := decrypt(key.Value)
			if decErr != nil {
				return nil, err
			}
			results[i].Value = dec
		}
	}
	return results, nil
}
