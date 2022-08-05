/*
  Onix Config Manager - Mongo Library
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	ErrDocumentAlreadyExists = errors.New("mongo: the document already exists")
	ErrDocumentNotFound      = errors.New("mongo: the document was not found")
)

// Database manage MongoDb connections
type Database struct {
	options *options.ClientOptions
}

func New(connString string) *Database {
	return &Database{
		options: options.Client().ApplyURI(connString),
	}
}

// ctx create a context with timeout of 30 seconds
func ctx() (context.Context, context.CancelFunc) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	return c, cancel
}

func (db *Database) Count(dbName string, collection Collection, filter interface{}) (int64, error) {
	c, cancel := ctx()
	defer cancel()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return 0, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	return coll.CountDocuments(c, filter)
}

// InsertObject insert a nameable object in the specified collection
func (db *Database) InsertObject(dbName string, collection Collection, obj Nameable) (interface{}, error) {
	item, err := db.FindByName(dbName, collection, obj.GetName())
	if err != nil {
		return nil, err
	}
	// if the key was found
	if item.Err() != mongo.ErrNoDocuments {
		return nil, ErrDocumentAlreadyExists
	}
	c, _ := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Printf(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	// insert the key
	result, insertErr := coll.InsertOne(c, obj)
	if insertErr != nil {
		return nil, fmt.Errorf("cannot insert object into %s collection: %s", collection, err)
	}
	return result.InsertedID, nil
}

func (db *Database) UpdateObject(dbName string, collection Collection, obj Nameable) (interface{}, error) {
	item, err := db.FindByName(dbName, collection, obj.GetName())
	if err != nil {
		return nil, err
	}
	// if the key was not found
	if item.Err() == mongo.ErrNoDocuments {
		return nil, ErrDocumentNotFound
	}
	c, _ := ctx()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	result, updateErr := coll.UpdateOne(c, bson.M{"_id": obj.GetName()}, bson.M{"$set": obj})
	if updateErr != nil {
		return nil, updateErr
	}
	return result, nil
}

func (db *Database) UpsertObject(dbName string, collection Collection, obj Nameable) (result interface{}, err error, resultCode int) {
	result, err = db.InsertObject(dbName, collection, obj)
	if err == nil {
		resultCode = 201
	} else if err == ErrDocumentAlreadyExists {
		// performs an update instead
		result, err = db.UpdateObject(dbName, collection, obj)
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

func (db *Database) DeleteByName(dbName string, collection Collection, name string) (interface{}, error) {
	c, cancel := ctx()
	defer cancel()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	return coll.DeleteOne(c, bson.M{"_id": name})
}

// FindByName find an object by name
func (db *Database) FindByName(dbName string, collection Collection, name string) (*mongo.SingleResult, error) {
	c, cancel := ctx()
	defer cancel()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	item := coll.FindOne(c, bson.M{"_id": name})
	return item, nil
}

// FindOne object using the specified filter, examples of filter expressions below"
// filtering using  a slice
//   filter := bson.D{{"attribute_name1", "attribute_value1"}, {"attribute_name2", "attribute_value2"}, ...}
// filtering using a map
//   filter := bson.M{"attribute_name1": "attribute_value1", "attribute_name2": "attribute_value2", ...}
func (db *Database) FindOne(dbName string, collection Collection, filter interface{}) (*mongo.SingleResult, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is required")
	}
	client, err := mongo.Connect(context.Background(), db.options)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, context.Background())
	coll := client.Database(dbName).Collection(string(collection))
	result := coll.FindOne(context.Background(), filter)
	return result, nil
}

// FindMany find a number of objects matching the specified filter
func (db *Database) FindMany(dbName string, collection Collection, filter bson.M, query Query) error {
	if filter == nil {
		filter = bson.M{}
	}
	client, err := mongo.Connect(context.Background(), db.options)
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, context.Background())
	coll := client.Database(dbName).Collection(string(collection))
	cursor, findErr := coll.Find(context.Background(), filter)
	defer func() {
		// only closes the cursor if it exists
		if cursor != nil {
			err = cursor.Close(context.Background())
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
	if findErr != nil {
		return findErr
	}
	return query(cursor)
}

type Query func(cursor *mongo.Cursor) error

func (db *Database) FindAll(dbName string, collection Collection, results interface{}) error {
	c, cancel := ctx()
	defer cancel()
	client, err := mongo.Connect(c, db.options)
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(client, c)
	coll := client.Database(dbName).Collection(string(collection))
	cursor, err := coll.Find(c, make(map[string]interface{}))
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}(cursor, c)
	if err != nil {
		return err
	}
	return cursor.All(c, &results)
}

// ObjectExists checks if an object exists in the specified collection
func (db *Database) ObjectExists(dbName string, collection Collection, name string) bool {
	item, err := db.FindByName(dbName, collection, name)
	if err != nil {
		log.Printf("cannot retrieve item %s in collection %s: %s\n", name, collection, err)
		return false
	}
	return item.Err() != mongo.ErrNoDocuments
}
