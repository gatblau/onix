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
func processAlerts(data template.Alerts, get getItem, put putItem, partition string) error {
	// sort the incoming alerts by StartsAt time
	alerts := NewTimeSortedAlerts(data)
	sort.Sort(alerts)
	log.Debug().Msgf("* alerts have been sorted")
	// creates a map to keep track of existing items and values
	items := make(map[string]*oxc.Item)
	// loops through tha alerts
	for _, alert := range alerts {
		// extract values from the alert
		v, err := values(alert)
		// if the alert did not have all required information
		if err != nil {
			// stops any processing
			return err
		}
		// write debug info about successful data extraction
		log.Debug().Msgf("* extracted values for alert with fingerprint '%s'", alert.Fingerprint)
		log.Debug().Msgf("* platform value = '%s'", v["platform"])
		log.Debug().Msgf("* service value = '%s'", v["service"])
		log.Debug().Msgf("* facet value = '%s'", v["facet"])

		location := v["location"]
		// if the alert does not have a specific location
		if len(location) == 0 {
			// fill the blank
			location = "_"
		}
		// build the natural key for the service item
		serviceKey := key(v["platform"], v["service"], v["facet"], location)
		log.Debug().Msgf("* service key '%s' created", serviceKey)

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
			if err != nil {
				log.Debug().Msgf("* fail to fetch information for item with key '%s': '%s'", serviceKey, err)
			}
		}
		var startsAt time.Time
		var shouldRegisterEvent bool

		// if there is a pre-existing item
		if serviceItem != nil {
			// extract the event time
			if startsAtStr, ok := serviceItem.Attribute["time"].(string); ok {
				startsAt, err = time.Parse(time.RFC3339Nano, startsAtStr)
				if err != nil {
					log.Warn().Msgf("failed to parse startsAt time '%s': %s", startsAtStr, err)
				}
			} else if startsAt, ok = serviceItem.Attribute["time"].(time.Time); !ok {
				// cannot parse time!
				log.Warn().Msgf("failed to parse startsAt time for event %s", serviceKey)
			}
		}

		// Decides if the event should be registered in Onix using the logic below:
		// 1) there is not a record of the event in Onix; or
		// 2) the event registered in Onix has a status that is different from the status in the received alert; and
		// 3) the received alert has occurred after the last event registered in Onix
		shouldRegisterEvent = (serviceItem == nil || (serviceItem != nil && serviceItem.Attribute["status"].(string) != v["status"])) && startsAt.Before(alert.StartsAt)
		log.Debug().Msgf("* event for item '%s' should be registered?: %t", serviceKey, shouldRegisterEvent)

		// if the item already recorded occurred before the current alert
		// or the no item recorded yet (startsAt = beginning of time)
		if shouldRegisterEvent {
			// then record the new alert
			log.Info().Msgf("recording event => %s:%s:%s:%s", v["service"], v["facet"], location, v["status"])
			serviceItem = &oxc.Item{
				Key:         serviceKey,
				Name:        fmt.Sprintf("%s Service", v["service"]),
				Description: v["description"],
				Type:        SeSServiceItemType,
				Attribute: map[string]interface{}{
					"service":     v["service"],
					"status":      v["status"],
					"description": v["description"],
					"location":    v["location"],
					"time":        alert.StartsAt,
					"platform":    v["platform"],
					"facet":       v["facet"],
				},
				Partition: partition,
			}
			result, err := put(serviceItem)
			// do we have a business error?
			if result != nil && result.Error {
				return errors.New(fmt.Sprintf("Onix http put failed: %s", result.Message))
			}
			if err != nil {
				return errors.New(fmt.Sprintf("Onix http put failed: %s", err))
			}
			// update the internal cache
			items[serviceKey] = serviceItem
		} else {
			log.Trace().Msgf("discarding event => %s:%s:%s:%s", v["service"], v["facet"], location, v["status"])
		}
	}
	return nil
}

// extract values from the alert
func values(alert template.Alert) (values map[string]string, err error) {
	var result = make(map[string]string)

	err, result["platform"] = kValue(alert.Labels, "platform")
	if err != nil {
		return result, err
	}
	err, result["service"] = kValue(alert.Labels, "service")
	if err != nil {
		return result, err
	}
	err, result["status"] = kValue(alert.Labels, "status")
	if err != nil {
		return result, err
	}
	err, result["description"] = kValue(alert.Labels, "description")
	if err != nil {
		return result, err
	}
	err, result["facet"] = kValue(alert.Labels, "facet")
	if err != nil {
		return result, err
	}
	// add any annotations
	for key, value := range alert.Annotations {
		// if the annotation key is not a label (e.g. not in the result map already)
		if result[key] == "" {
			// add the annotation value
			result[key] = value
		} else {
			// skip and issue a warning
			log.Warn().Msgf("skipping annotation '%s' as a label with such name was found", key)
		}
	}
	return result, nil
}
