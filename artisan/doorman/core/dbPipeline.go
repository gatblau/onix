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
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func FindPipeline(pipeName string) (*types.Pipeline, error) {
	db := NewDb()
	result, err := db.FindByName(types.PipelineCollection, pipeName)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pipeline %s: %s", pipeName, err)
	}
	pipeConf := new(types.PipelineConf)
	err = result.Decode(pipeConf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode pipeline %s: %s", pipeName, err)
	}
	result, err = db.FindByName(types.InRouteCollection, pipeConf.InboundRoute)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve inbound route %s: %s", pipeConf.InboundRoute, err)
	}
	inRoute := new(types.InRoute)
	err = result.Decode(inRoute)
	if err != nil {
		return nil, fmt.Errorf("cannot decode inbound route %s: %s", pipeConf.InboundRoute, err)
	}
	result, err = db.FindByName(types.OutRouteCollection, pipeConf.OutboundRoute)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve outbound route %s: %s", pipeConf.OutboundRoute, err)
	}
	outRoute := new(types.OutRoute)
	err = result.Decode(outRoute)
	if err != nil {
		return nil, fmt.Errorf("cannot decode outbound route %s: %s", pipeConf.OutboundRoute, err)
	}
	pipe := &types.Pipeline{
		Name:          pipeConf.Name,
		InboundRoute:  *inRoute,
		OutboundRoute: *outRoute,
		Commands:      pipeConf.Commands,
	}
	return pipe, nil
}

func UpsertPipeline(pipe types.PipelineConf) (error, int) {
	db := NewDb()
	_, err := db.FindByName(types.InRouteCollection, pipe.InboundRoute)
	if err != nil {
		return fmt.Errorf("cannot find inbound route %s for pipeline %s: %s", pipe.InboundRoute, pipe.Name, err), http.StatusBadRequest
	}
	_, err = db.FindByName(types.OutRouteCollection, pipe.OutboundRoute)
	if err != nil {
		return fmt.Errorf("cannot find outbound route %s for pipeline %s: %s", pipe.OutboundRoute, pipe.Name, err), http.StatusBadRequest
	}
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.PipelineCollection, pipe)
	if err != nil {
		return fmt.Errorf("cannot update pipeline in database"), resultCode
	}
	return nil, resultCode
}

func FindAllPipelines() ([]types.PipelineConf, error) {
	db := NewDb()
	var pipelines []types.PipelineConf
	if err := db.FindMany(types.PipelineCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &pipelines)
	}); err != nil {
		return nil, err
	}
	return pipelines, nil
}
