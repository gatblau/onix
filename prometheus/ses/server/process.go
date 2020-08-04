/*
   Onix Config Manager - SeS - Onix Webhook Receiver for AlertManager
   Copyright (c) 2020 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
package server

import (
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/prometheus/alertmanager/template"
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

// these functions decouple key backend resources to enable unit testing
// getItem: gets an item from Onix based on its key
type getItem func(item *oxc.Item) (*oxc.Item, error)

// putItem: puts an item in Onix
type putItem func(item *oxc.Item) (*oxc.Result, error)

// process the received alerts
func processAlerts(data template.Alerts, get getItem, put putItem) error {
	// sort the incoming alerts by StartsAt time
	alerts := NewTimeSortedAlerts(data)
	sort.Sort(alerts)
	// creates a map to keep track of existing items and values
	items := make(map[string]*oxc.Item)
	// loops through tha alerts
	for _, alert := range alerts {
		// extract values from the alert
		platform, service, status, description, uri, err := values(alert)
		// build the natural key for the service item
		serviceKey := key(platform, service, uri)
		// check if there is a previous record of the service item in the tracking map
		serviceItem := items[serviceKey]
		// if the item is not in the map
		if serviceItem == nil {
			// try and fetch it from Onix
			serviceItem, err = get(&oxc.Item{Key: serviceKey})
			// if there is an item
			if err == nil && serviceItem != nil {
				// add it to the tracking map
				items[serviceKey] = serviceItem
			}
		}
		var startsAt time.Time
		var shouldRegisterEvent bool

		// if there is a pre-existing item
		if serviceItem != nil {
			// extract the event time
			if startsAtStr, ok := serviceItem.Attribute["time"].(string); ok {
				startsAt, err = time.Parse(time.RFC3339Nano, startsAtStr)
			} else if startsAt, ok = serviceItem.Attribute["time"].(time.Time); !ok {
				// cannot parse time!
				log.Warn().Msgf("failed to parse startsAt time for event %s", serviceKey)
			}
			if err != nil {
				// discard the startsAt time
				log.Warn().Msgf("failed to parse startsAt time: %s", err)
			}
		}

		// Decides if the event should be registered in Onix using the logic below:
		// 1) there is not a record of the event in Onix; or
		// 2) the event registered in Onix has a status that is different from the status in the received alert; and
		// 3) the received alert has occurred after the last event registered in Onix
		shouldRegisterEvent = (serviceItem == nil || (serviceItem != nil && serviceItem.Attribute["status"].(string) != status)) && startsAt.Before(alert.StartsAt)

		// if the item already recorded occurred before the current alert
		// or the no item recorded yet (startsAt = beginning of time)
		if shouldRegisterEvent {
			// then record the new alert
			log.Info().Msgf("recording event: service %s:%s is %s", service, uri, status)
			serviceItem = &oxc.Item{
				Key:         serviceKey,
				Name:        fmt.Sprintf("%s Service", service),
				Description: description,
				Type:        SeSServiceItemType,
				Attribute: map[string]interface{}{
					"name":        service,
					"status":      status,
					"description": description,
					"uri":         uri,
					"time":        alert.StartsAt,
					"platform":    platform,
				},
			}
			result, err := put(serviceItem)
			if err != nil {
				return err
			}
			if result.Error {
				return errors.New(fmt.Sprintf("cannot update service status: %s", result.Message))
			}
			// update the internal cache
			items[serviceKey] = serviceItem
		} else {
			log.Trace().Msgf("discarding event: service %s:%s is %s", service, uri, status)
		}
	}
	return nil
}

// extract values from the alert
func values(alert template.Alert) (platform string, service string, status string, description string, uri string, err error) {
	var ok bool
	ok, platform = kValue(alert.Labels, "platform")
	if !ok {
		return platform, "", "", "", "", errors.New(fmt.Sprintf("cannot find 'platform' annotation in alert '%s'", alert))
	}
	ok, service = kValue(alert.Labels, "service")
	if !ok {
		return platform, service, "", "", "", errors.New(fmt.Sprintf("cannot find 'service' annotation in alert '%s'", alert))
	}
	ok, status = kValue(alert.Labels, "status")
	if !ok {
		return platform, "", "", "", "", errors.New(fmt.Sprintf("cannot find 'status' annotation in alert '%s'", alert))
	}
	ok, uri = kValue(alert.Labels, "uri")
	if !ok {
		return "", "", "", "", "", errors.New(fmt.Sprintf("cannot find 'uri' annotation in alert '%s'", alert))
	}
	ok, description = kValue(alert.Labels, "description")
	if !ok {
		return "", "", "", "", "", errors.New(fmt.Sprintf("cannot find 'description' annotation in alert '%s'", alert))
	}
	return platform, service, status, description, uri, nil
}
