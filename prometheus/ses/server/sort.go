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
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"strings"
)

// type to sort alert by StartsAt time (implement sort interface)
type TimeSortedAlerts []template.Alert

func NewTimeSortedAlerts(alerts []template.Alert) TimeSortedAlerts {
	result := make(TimeSortedAlerts, 0)
	for _, alert := range alerts {
		result = append(result, alert)
	}
	return result
}

func (a TimeSortedAlerts) Len() int {
	return len(a)
}

func (a TimeSortedAlerts) Less(i, j int) bool {
	return a[i].StartsAt.Before(a[j].StartsAt)
}

func (a TimeSortedAlerts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// check if the passed-in map contains the specified key and returns its value
func kValue(values template.KV, key string) (bool, string) {
	for k, v := range values {
		if strings.ToLower(key) == strings.ToLower(k) {
			return true, v
		}
	}
	return false, ""
}

// get the service unique natural key
func key(platform string, service string, facet string, location string) string {
	return fmt.Sprintf("%s_%s_%s_%s", platform, service, facet, strings.Replace(strings.Replace(location, ":", "_", -1), ".", "_", -1))
}
