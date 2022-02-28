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
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (db *Db) FindPipeline(pipeName string) (*types.Pipeline, error) {
	result, err := db.FindByName(types.PipelineCollection, pipeName)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pipeline %s: %s", pipeName, err)
	}
	pipeConf := new(types.PipelineConf)
	err = result.Decode(pipeConf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode pipeline %s: %s", pipeName, err)
	}
	var inRoutes []types.InRoute
	for _, route := range pipeConf.InboundRoutes {
		result, err = db.FindByName(types.InRouteCollection, route)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve inbound route %s: %s", route, err)
		}
		inRoute := new(types.InRoute)
		err = result.Decode(inRoute)
		if err != nil {
			return nil, fmt.Errorf("cannot decode inbound route %s: %s", route, err)
		}
		inRoutes = append(inRoutes, *inRoute)
	}
	var outRoutes []types.OutRoute
	for _, route := range pipeConf.OutboundRoutes {
		result, err = db.FindByName(types.OutRouteCollection, route)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve outbound route %s: %s", route, err)
		}
		outRoute := new(types.OutRoute)
		err = result.Decode(outRoute)
		if err != nil {
			return nil, fmt.Errorf("cannot decode outbound route %s: %s", route, err)
		}
		outRoutes = append(outRoutes, *outRoute)
	}
	var cmds []types.Command
	for _, cmd := range pipeConf.Commands {
		result, err = db.FindByName(types.CommandsCollection, cmd)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve command %s: %s", cmd, err)
		}
		cmdObj := new(types.Command)
		err = result.Decode(cmdObj)
		if err != nil {
			return nil, fmt.Errorf("cannot decode command %s: %s", cmd, err)
		}
		cmds = append(cmds, *cmdObj)
	}
	successN, err := db.FindNotification(pipeConf.SuccessNotification)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve success notification %s: %s", pipeConf.SuccessNotification, err)
	}
	cmdFailedN, err := db.FindNotification(pipeConf.CmdFailedNotification)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve command failed notification %s: %s", pipeConf.CmdFailedNotification, err)
	}
	errorN, err := db.FindNotification(pipeConf.ErrorNotification)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve error notification %s: %s", pipeConf.ErrorNotification, err)
	}
	pipe := &types.Pipeline{
		Name:                  pipeConf.Name,
		InboundRoutes:         inRoutes,
		OutboundRoutes:        outRoutes,
		Commands:              cmds,
		SuccessNotification:   successN,
		CmdFailedNotification: cmdFailedN,
		ErrorNotification:     errorN,
	}
	return pipe, nil
}

func (db *Db) UpsertPipeline(pipe types.PipelineConf) (error, int) {
	var err error
	for _, route := range pipe.InboundRoutes {
		_, err = db.FindByName(types.InRouteCollection, route)
		if err != nil {
			return fmt.Errorf("cannot find inbound route %s for pipeline %s: %s", route, pipe.Name, err), http.StatusBadRequest
		}
	}
	for _, route := range pipe.OutboundRoutes {
		_, err = db.FindByName(types.OutRouteCollection, route)
		if err != nil {
			return fmt.Errorf("cannot find outbound route %s for pipeline %s: %s", route, pipe.Name, err), http.StatusBadRequest
		}
	}
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.PipelineCollection, pipe)
	if err != nil {
		return fmt.Errorf("cannot update pipeline in database"), resultCode
	}
	return nil, resultCode
}

func (db *Db) FindAllPipelines() ([]types.PipelineConf, error) {
	var pipelines []types.PipelineConf
	if err := db.FindMany(types.PipelineCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &pipelines)
	}); err != nil {
		return nil, err
	}
	return pipelines, nil
}

func (db *Db) FindPipelinesByInboundURI(uri string) ([]types.Pipeline, error) {
	var (
		pipes    []types.Pipeline
		routes   []types.InRoute
		pipeline *types.Pipeline
		err      error
	)
	routes, err = db.FindInboundRoutesByURI(uri)
	if err != nil {
		return nil, err
	}
	var pipeConfs []types.PipelineConf
	for _, route := range routes {
		// any pipeline having route.Name in their inbound routes array
		filter := bson.M{"inbound_routes": bson.M{"$all": []string{route.Name}}}
		if err = db.FindMany(types.PipelineCollection, filter, func(cursor *mongo.Cursor) error {
			return cursor.All(context.Background(), &pipeConfs)
		}); err != nil {
			return nil, err
		}
	}
	for _, conf := range pipeConfs {
		pipeline, err = db.FindPipeline(conf.Name)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, *pipeline)
	}
	return pipes, nil
}

func (db *Db) MatchPipelines(serviceId, bucketName string) ([]types.Pipeline, error) {
	var (
		pipes    []types.Pipeline
		routes   []types.InRoute
		pipeline *types.Pipeline
		err      error
	)
	routes, err = db.MatchInboundRoutes(serviceId, bucketName)
	if err != nil {
		return nil, err
	}
	var pipeConfs []types.PipelineConf
	for _, route := range routes {
		// any pipeline having route.Name in their inbound routes array
		filter := bson.M{"inbound_routes": bson.M{"$all": []string{route.Name}}}
		if err = db.FindMany(types.PipelineCollection, filter, func(cursor *mongo.Cursor) error {
			return cursor.All(context.Background(), &pipeConfs)
		}); err != nil {
			return nil, err
		}
	}
	for _, conf := range pipeConfs {
		pipeline, err = db.FindPipeline(conf.Name)
		if err != nil {
			return nil, err
		}
		if err = pipeline.Valid(); err != nil {
			return nil, err
		}
		pipes = append(pipes, *pipeline)
	}
	return pipes, nil
}
