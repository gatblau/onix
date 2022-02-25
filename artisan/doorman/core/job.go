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
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (db *Db) StartJob(pipeline *types.Pipeline, processor *Processor) (string, *time.Time, error) {
	jobNo := uuid.New().String()
	startTime := time.Now().UTC()
	_, err := db.InsertObject(types.JobsCollection, &types.Job{
		Number:    jobNo,
		ServiceId: processor.serviceId,
		Bucket:    processor.bucketName,
		Folder:    processor.folderName,
		Pipeline:  pipeline,
		Status:    "started",
		Started:   &startTime,
	})
	if err != nil {
		return "", nil, err
	}
	return jobNo, &startTime, nil
}

func (db *Db) CompleteJob(started *time.Time, pipeline *types.Pipeline, processor *Processor) error {
	completedTime := time.Now().UTC()
	_, err, _ := db.UpsertObject(types.JobsCollection, &types.Job{
		Number:    processor.jobNo,
		ServiceId: processor.serviceId,
		Bucket:    processor.bucketName,
		Folder:    processor.folderName,
		Status:    "completed",
		Pipeline:  pipeline,
		Log:       processor.logs(),
		Started:   started,
		Completed: &completedTime,
	})
	return err
}

func (db *Db) FailJob(started *time.Time, pipeline *types.Pipeline, processor *Processor) error {
	completedTime := time.Now().UTC()
	_, err, _ := db.UpsertObject(types.JobsCollection, &types.Job{
		Number:    processor.jobNo,
		ServiceId: processor.serviceId,
		Bucket:    processor.bucketName,
		Folder:    processor.folderName,
		Status:    "failed",
		Pipeline:  pipeline,
		Log:       processor.logs(),
		Started:   started,
		Completed: &completedTime,
	})
	return err
}

func (db *Db) FindTopJobs(count int) ([]types.Job, error) {
	var jobs []types.Job
	if err := db.FindMany(types.JobsCollection, nil, func(cursor *mongo.Cursor) error {
		return cursor.All(context.TODO(), &jobs)
	}); err != nil {
		return nil, err
	}
	return jobs, nil
}
