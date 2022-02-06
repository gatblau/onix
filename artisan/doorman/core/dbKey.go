/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Keys return the keys collection
func (db *Db) Keys(client *mongo.Client) *mongo.Collection {
	return client.Database(DbName).Collection("keys")
}

func (db *Db) KeyExists(name string) bool {
	c := ctx()
	// connect to mongo
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		log.Printf("cannot connect to the database: %s", err)
	}
	defer client.Disconnect(c)
	// check if a key with the same name already exists
	item := db.Keys(client).FindOne(c, bson.M{"_id": name})
	// if the key was found
	return item.Err() != mongo.ErrNoDocuments
}

func (db *Db) NewKey(key *types.Key) (interface{}, error) {
	// if the key is private
	if key.IsPrivate {
		// encrypt it
		enc, encErr := encrypt(key.Value)
		if encErr != nil {
			return nil, encErr
		}
		key.Value = enc
	}
	c := ctx()
	// connect to mongo
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(c)
	// check if a key with the same name already exists
	item := db.Keys(client).FindOne(c, bson.M{"name": key.Name})
	// if the key was found
	if item.Err() != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("key with name %s already exist", key.Name)
	}
	// insert the key
	result, insertErr := db.Keys(client).InsertOne(c, key)
	if insertErr != nil {
		return nil, fmt.Errorf("cannot insert key: %s", err)
	}
	return result.InsertedID, nil
}

func (db *Db) FindKeys(filter bson.M) ([]types.Key, error) {
	c := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(c)
	cursor, findErr := db.Keys(client).Find(ctx(), filter)
	if findErr != nil {
		return nil, findErr
	}
	var results []types.Key
	// return elements that match the criteria
	if err = cursor.All(ctx(), &results); err != nil {
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
