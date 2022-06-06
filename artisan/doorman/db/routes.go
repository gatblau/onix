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
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *Database) FindInboundRoutesByURI(uri string) ([]types.InRoute, error) {
	var routes []types.InRoute
	if err := db.FindMany(types.InRouteCollection, bson.M{"bucket_uri": uri}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}

func (db *Database) MatchInboundRoutes(serviceId, bucketName string) ([]types.InRoute, error) {
	var routes []types.InRoute
	// first try to mach routes with any bucket (*)
	err := db.FindMany(types.InRouteCollection, bson.M{"service_id": serviceId, "bucket_name": "*"}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	})
	if err != nil {
		return nil, fmt.Errorf("cannot find route with wildcard bucket name: %s\n", err)
	}
	if routes != nil {
		if len(routes) > 1 {
			return []types.InRoute{routes[0]}, nil
		} else {
			return routes, nil
		}
	} else { // if no match is found, then run a search using bucket name
		if err = db.FindMany(types.InRouteCollection, bson.M{"service_id": serviceId, "bucket_name": bucketName}, func(cursor *mongo.Cursor) error {
			return cursor.All(context.Background(), &routes)
		}); err != nil {
			return nil, err
		}
	}
	return routes, nil
}

func (db *Database) FindInboundRoutesByWebHookToken(token string) ([]types.InRoute, error) {
	var routes []types.InRoute
	if err := db.FindMany(types.InRouteCollection, bson.M{"webhook_token": token}, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}

func (db *Database) FindAllInRoutes() ([]types.InRoute, error) {
	var routes []types.InRoute
	if err := db.FindMany(types.InRouteCollection, nil, func(cursor *mongo.Cursor) error {
		return cursor.All(context.Background(), &routes)
	}); err != nil {
		return nil, err
	}
	return routes, nil
}
