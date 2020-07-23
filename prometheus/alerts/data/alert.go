//
//    Onix Config Manager - Alert Manager Webhook Receiver for Onix
//    Copyright (c) 2018-2020 by www.gatblau.org
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//    Unless required by applicable law or agreed to in writing, software distributed under
//    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//    either express or implied.
//    See the License for the specific language governing permissions and limitations under the License.
//
//    Contributors to this project, hereby assign copyright in this code to the project,
//    to be licensed under the same terms as the rest of the code.
//
package data

import "time"

// the payload in a webhook call
type AlertGroup struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`

	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alerts is a slice of Alert
type Alerts []Alert

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
}
