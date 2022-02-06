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
	"go.mongodb.org/mongo-driver/mongo/options"
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
