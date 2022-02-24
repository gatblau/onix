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
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindKeys retrieves one or more keys matching the specifies criteria decrypting the value of any private key
func (db *Db) FindKeys(filter bson.M) ([]types.Key, error) {
	var keys []types.Key
	if err := db.FindMany(types.KeysCollection, filter, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &keys)
	}); err != nil {
		return nil, err
	}
	for i, key := range keys {
		if key.IsPrivate {
			dec, decErr := decrypt(key.Value)
			if decErr != nil {
				return nil, decErr
			}
			keys[i].Value = dec
		}
	}
	return keys, nil
}

// FindKeyByName retrieves the key with the specified name decrypting its value if it is a private key
func (db *Db) FindKeyByName(name string) (*types.Key, error) {
	var key types.Key
	result, err := db.FindByName(types.KeysCollection, name)
	if err != nil {
		return nil, err
	}
	err = result.Decode(&key)
	if err != nil {
		return nil, err
	}
	if key.IsPrivate {
		dec, decErr := decrypt(key.Value)
		if decErr != nil {
			return nil, decErr
		}
		key.Value = dec
	}
	return &key, nil
}

func (db *Db) UpsertKey(key *types.Key) (error, int) {
	if key.IsPrivate {
		enc, err := encrypt(key.Value)
		if err != nil {
			return err, -1
		}
		key.Value = enc
	}
	_, err, resultCode := db.UpsertObject(types.KeysCollection, key)
	if err != nil {
		return err, resultCode
	}
	return nil, resultCode
}
