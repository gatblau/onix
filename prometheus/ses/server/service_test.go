package server

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"io/ioutil"
	"os/user"
	"sort"
	"testing"
)

func TestSeS_processAlerts(t *testing.T) {
	// get the current user
	usr, err := user.Current()
	if err != nil {
		t.Error(err)
		return
	}
	// load alerts from file
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/go/src/github.com/gatblau/onix/prometheus/ses/test/payload_up.json", usr.HomeDir))
	if err != nil {
		t.Error(err)
		return
	}
	// unmarshal alerts
	var payload template.Data
	json.Unmarshal([]byte(dat), &payload)
	if err != nil {
		t.Error(err)
		return
	}
	alerts := NewTimeSortedAlerts(payload.Alerts)
	for _, alert := range alerts {
		if alert.Annotations["service"] != "" {
			fmt.Printf("%s: %s\n", alert.Annotations["service"], alert.StartsAt)
		}
	}
	// now sort the alerts
	sort.Sort(alerts)
	for _, alert := range alerts {
		if alert.Annotations["service"] != "" {
			fmt.Printf("%s: %s\n", alert.Annotations["service"], alert.StartsAt)
		}
	}
}

const payload = `
{
  "receiver": "onix-webhook",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "EtcdIsUp",
        "instance": "127.0.0.1:22379",
        "job": "etcd",
        "severity": "info"
      },
      "annotations": {
        "description": "The etcd server is up",
        "instance": "127.0.0.1:22379",
        "service": "etcd",
        "status": "up"
      },
      "startsAt": "2020-07-28T14:23:22.095063475Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://localhost:9090/graph?g0.expr=up%7Bjob%3D%22etcd%22%7D+%3D%3D+1\u0026g0.tab=1",
      "fingerprint": "329292b9642baaba"
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "EtcdIsUp",
        "instance": "127.0.0.1:2379",
        "job": "etcd",
        "severity": "info"
      },
      "annotations": {
        "description": "The etcd server is up",
        "instance": "127.0.0.1:2379",
        "service": "etcd",
        "status": "up"
      },
      "startsAt": "2020-07-28T14:23:22.095063475Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://localhost:9090/graph?g0.expr=up%7Bjob%3D%22etcd%22%7D+%3D%3D+1\u0026g0.tab=1",
      "fingerprint": "f6b268ca5edb35f4"
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "EtcdIsUp",
        "instance": "127.0.0.1:32379",
        "job": "etcd",
        "severity": "info"
      },
      "annotations": {
        "description": "The etcd server is up",
        "instance": "127.0.0.1:32379",
        "service": "etcd",
        "status": "up"
      },
      "startsAt": "2020-07-28T14:23:22.095063475Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://localhost:9090/graph?g0.expr=up%7Bjob%3D%22etcd%22%7D+%3D%3D+1\u0026g0.tab=1",
      "fingerprint": "0398ec31c84fbe9d"
    }
  ],
  "groupLabels": {
    "job": "etcd"
  },
  "commonLabels": {
    "alertname": "EtcdIsUp",
    "job": "etcd",
    "severity": "info"
  },
  "commonAnnotations": {
    "description": "The etcd server is up",
    "service": "etcd",
    "status": "up"
  },
  "externalURL": "http://localhost:9093",
  "version": "4",
  "groupKey": "{}/{alertname=\"EtcdIsUp\"}:{job=\"etcd\"}",
  "truncatedAlerts": 0
}
`
