/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import (
	"time"
)

type Job struct {
	Number    string     `bson:"_id" json:"number" yaml:"number"`
	ServiceId string     `bson:"service_id" json:"service_id" yaml:"service_id"`
	Bucket    string     `bson:"bucket" json:"bucket" yaml:"bucket"`
	Folder    string     `bson:"folder" json:"folder" yaml:"folder"`
	Pipeline  *Pipeline  `bson:"pipeline" json:"pipeline" yaml:"pipeline"`
	Status    string     `bson:"status" json:"status" yaml:"status"`
	Log       []string   `bson:"log" json:"log" yaml:"log"`
	Started   *time.Time `bson:"started" json:"started" yaml:"started"`
	Completed *time.Time `bson:"completed" json:"completed" yaml:"completed"`
}

func (j *Job) GetName() string {
	return j.Number
}
