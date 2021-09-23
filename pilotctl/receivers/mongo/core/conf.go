package core

/*
  Onix Config Manager - MongoDb event receiver for Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import "os"

// getDbConnString get the connection string to the MongoDb database
// e.g. mongodb://localhost:27017
// e.g. mongodb://user:password@127.0.0.1:27017/dbname?keepAlive=true&poolSize=30&autoReconnect=true&socketTimeoutMS=360000&connectTimeoutMS=360000
func getDbConnString() string {
	value := os.Getenv("OX_MONGO_EVR_CONN")
	if len(value) == 0 {
		panic("OX_MONGO_EVR_CONN not defined")
	}
	return value
}
