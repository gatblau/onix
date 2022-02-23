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

func FindInboundRoutesByURI(uri string) ([]types.InRoute, error) {
	var routes []types.InRoute
	db := NewDb()
	if err := db.FindMany(types.InRouteCollection, bson.M{"bucket_uri": uri}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}

func FindInboundRoutesById(id string) ([]types.InRoute, error) {
	var routes []types.InRoute
	db := NewDb()
	if err := db.FindMany(types.InRouteCollection, bson.M{"bucket_id": id}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}

func FindInboundRoutesByWebHookToken(token string) ([]types.InRoute, error) {
	var routes []types.InRoute
	db := NewDb()
	if err := db.FindMany(types.InRouteCollection, bson.M{"webhook_token": token}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}

func FindAllInRoutes() ([]types.InRoute, error) {
	var routes []types.InRoute
	db := NewDb()
	if err := db.FindMany(types.InRouteCollection, nil, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}
