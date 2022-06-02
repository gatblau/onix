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
	"github.com/gatblau/onix/artisan/doorman/db"
	"github.com/gatblau/onix/artisan/doorman/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func StartJob(pipeline *types.Pipeline, process *Process) (string, *time.Time, error) {
	jobNo := uuid.New().String()
	startTime := time.Now().UTC()
	_, err := process.db.InsertObject(types.JobsCollection, &types.Job{
		Number:    jobNo,
		ServiceId: process.serviceId,
		Bucket:    process.bucketName,
		Folder:    process.folderName,
		Pipeline:  pipeline,
		Status:    "started",
		Started:   &startTime,
	})
	if err != nil {
		return "", nil, err
	}
	return jobNo, &startTime, nil
}

func CompleteJob(started *time.Time, pipeline *types.Pipeline, process *Process) error {
	completedTime := time.Now().UTC()
	_, err, _ := process.db.UpsertObject(types.JobsCollection, &types.Job{
		Number:    process.jobNo,
		ServiceId: process.serviceId,
		Bucket:    process.bucketName,
		Folder:    process.folderName,
		Status:    "completed",
		Pipeline:  pipeline,
		Log:       process.logs(),
		Started:   started,
		Completed: &completedTime,
	})
	return err
}

func FailJob(started *time.Time, pipeline *types.Pipeline, process *Process) error {
	completedTime := time.Now().UTC()
	_, err, _ := process.db.UpsertObject(types.JobsCollection, &types.Job{
		Number:    process.jobNo,
		ServiceId: process.serviceId,
		Bucket:    process.bucketName,
		Folder:    process.folderName,
		Status:    "failed",
		Pipeline:  pipeline,
		Log:       process.logs(),
		Started:   started,
		Completed: &completedTime,
	})
	return err
}

func FindTopJobs(count int, db *db.Database) ([]types.Job, error) {
	var jobs []types.Job
	if err := db.FindMany(types.JobsCollection, nil, func(cursor *mongo.Cursor) error {
		return cursor.All(context.TODO(), &jobs)
	}); err != nil {
		return nil, err
	}
	return jobs, nil
}
