/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package db

import (
	"context"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindKeys retrieves one or more keys matching the specifies criteria decrypting the value of any private key
func (db *Database) FindKeys(filter bson.M) ([]types.Key, error) {
	var keys []types.Key
	if err := db.FindMany(types.KeysCollection, filter, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &keys)
	}); err != nil {
		return nil, err
	}
	return keys, nil
}

// FindKeyByName retrieves the key with the specified name decrypting its value if it is a private key
func (db *Database) FindKeyByName(name string) (*types.Key, error) {
	var key types.Key
	result, err := db.FindByName(types.KeysCollection, name)
	if err != nil {
		return nil, err
	}
	err = result.Decode(&key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (db *Database) UpsertKey(key *types.Key) (error, int) {
	_, err, resultCode := db.UpsertObject(types.KeysCollection, key)
	if err != nil {
		return err, resultCode
	}
	return nil, resultCode
}
