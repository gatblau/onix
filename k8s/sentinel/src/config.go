/*
   Sentinel - Copyright (c) 2019 by www.gatblau.org

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
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

// Sentinel configuration
type Config struct {
	KubeConfig string
	LogLevel   string
	Publishers Publishers
	Observe    Observe
	Platform   string
}

// the configuration for the event publishers
type Publishers struct {
	Publisher string
	Logger    Logger
	Webhook   []Webhook
	Broker    Broker
}

// the configuration for the logger publisher
type Logger struct {
	OutputTo  string
	LogFolder string
}

// the configuration for the web hook publisher
type Webhook struct {
	URI                string
	Authentication     string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// the configuration for the message broker publisher
type Broker struct {
	Addr        string
	Brokers     string
	Verbose     bool
	Certificate string
	Key         string
	CA          string
	Verify      bool
}

// the type of objects that can be observed by the controller
type Observe struct {
	Service               bool
	Pod                   bool
	PersistentVolume      bool
	PersistentVolumeClaim bool
	Namespace             bool
	Deployment            bool
	ReplicationController bool
	ReplicaSet            bool
	DaemonSet             bool
	Job                   bool
	Secret                bool
	ConfigMap             bool
	Ingress               bool
	ServiceAccount        bool
	ClusterRole           bool
	ResourceQuota         bool
	NetworkPolicy         bool
}

// creates a new configuration file passed by value
// to avoid thread sync issues
func NewConfig() (Config, error) {
	logrus.Infof("Loading configuration.")
	v := viper.New()
	// loads the configuration file
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	err := v.ReadInConfig() // find and read the config file
	if err != nil {         // handle errors reading the config file
		logrus.Errorf("Fatal error config file: %s \n", err)
		return Config{}, err
	}

	// binds all environment variables to make it container friendly
	v.AutomaticEnv()
	v.SetEnvPrefix("SL")
	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = v.BindEnv("KubeConfig")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Platform")
	_ = v.BindEnv("Publishers.Publisher")
	_ = v.BindEnv("Publishers.Logger.OutputTo")
	_ = v.BindEnv("Publishers.Logger.LogFolder")
	_ = v.BindEnv("Publishers.Broker.Addr")
	_ = v.BindEnv("Publishers.Broker.Brokers")
	_ = v.BindEnv("Publishers.Broker.Verbose")
	_ = v.BindEnv("Publishers.Broker.Certificate")
	_ = v.BindEnv("Publishers.Broker.Name")
	_ = v.BindEnv("Publishers.Broker.CA")
	_ = v.BindEnv("Publishers.Broker.Verify")
	_ = v.BindEnv("Observe.Service")
	_ = v.BindEnv("Observe.Pod")
	_ = v.BindEnv("Observe.PersistentVolume")
	_ = v.BindEnv("Observe.PersistentVolumeClaim")
	_ = v.BindEnv("Observe.Namespace")
	_ = v.BindEnv("Observe.Deployment")
	_ = v.BindEnv("Observe.ReplicationController")
	_ = v.BindEnv("Observe.ReplicaSet")
	_ = v.BindEnv("Observe.DaemonSet")
	_ = v.BindEnv("Observe.Job")
	_ = v.BindEnv("Observe.AppSecret")
	_ = v.BindEnv("Observe.ConfigMap")
	_ = v.BindEnv("Observe.Ingress")
	_ = v.BindEnv("Observe.ServiceAccount")
	_ = v.BindEnv("Observe.ClusterRole")
	_ = v.BindEnv("Observe.ResourceQuota")
	_ = v.BindEnv("Observe.NetworkPolicy")

	// creates a config struct and populate it with values
	c := new(Config)

	// general configuration
	c.KubeConfig = v.GetString("KubeConfig")
	c.LogLevel = v.GetString("LogLevel")
	c.Platform = v.GetString("Platform")

	// publishers configuration
	c.Publishers.Publisher = v.GetString("Publishers.Publisher")

	// logger publisher configuration
	c.Publishers.Logger.OutputTo = v.GetString("Publishers.Logger.OutputTo")
	c.Publishers.Logger.LogFolder = v.GetString("Publishers.Logger.LogFolder")

	// webhook publisher configuration - loads array of tables in TOML
	// load all configured webhooks
	hooks := v.Get("Publishers.Webhook")
	if hooks != nil {
		whs := hooks.([]interface{})
		c.Publishers.Webhook = make([]Webhook, len(whs))
		for i := 0; i < len(whs); i++ {
			wh := whs[i].(map[string]interface{})
			h := Webhook{
				URI:                wh["URI"].(string),
				Username:           wh["Username"].(string),
				Password:           wh["Password"].(string),
				Authentication:     wh["Authentication"].(string),
				InsecureSkipVerify: wh["InsecureSkipVerify"].(bool),
			}
			// have to do ad-hoc binding of the array as Viper currently does not support
			// binding of TOML Array of Tables now try and bind any env variable following
			// the format PUBLISHERS_WEBHOOK_N_URI, etc where N is the array index
			value := os.Getenv(fmt.Sprintf("SL_PUBLISHERS_WEBHOOK_%s_URI", strconv.Itoa(i)))
			if len(value) > 0 {
				h.URI = value
			}
			value = os.Getenv(fmt.Sprintf("SL_PUBLISHERS_WEBHOOK_%s_USERNAME", strconv.Itoa(i)))
			if len(value) > 0 {
				h.Username = value
			}
			value = os.Getenv(fmt.Sprintf("SL_PUBLISHERS_WEBHOOK_%s_PASSWORD", strconv.Itoa(i)))
			if len(value) > 0 {
				h.Password = value
			}
			value = os.Getenv(fmt.Sprintf("SL_PUBLISHERS_WEBHOOK_%s_AUTHENTICATION", strconv.Itoa(i)))
			if len(value) > 0 {
				h.Authentication = value
			}
			value = os.Getenv(fmt.Sprintf("SL_PUBLISHERS_WEBHOOK_%s_INSECURESKIPVERIFY", strconv.Itoa(i)))
			if len(value) > 0 {
				h.InsecureSkipVerify, _ = strconv.ParseBool(value)
			}
			c.Publishers.Webhook[i] = h
		}
	}

	// broker publisher configuration
	c.Publishers.Broker.Addr = v.GetString("Publishers.Broker.Addr")
	c.Publishers.Broker.Brokers = v.GetString("Publishers.Broker.Brokers")
	c.Publishers.Broker.Certificate = v.GetString("Publishers.Broker.Certificate")
	c.Publishers.Broker.Key = v.GetString("Publishers.Broker.Name")
	c.Publishers.Broker.CA = v.GetString("Publishers.Broker.CA")
	c.Publishers.Broker.Verbose = v.GetBool("Publishers.Broker.Verbose")
	c.Publishers.Broker.Verify = v.GetBool("Publishers.Broker.Verify")

	// observable objects configuration
	c.Observe.Service = v.GetBool("Observe.Service")
	c.Observe.Pod = v.GetBool("Observe.Pod")
	c.Observe.PersistentVolume = v.GetBool("Observe.PersistentVolume")
	c.Observe.PersistentVolumeClaim = v.GetBool("Observe.PersistentVolumeClaim")
	c.Observe.Namespace = v.GetBool("Observe.Namespace")
	c.Observe.ConfigMap = v.GetBool("Observe.ConfigMap")
	c.Observe.DaemonSet = v.GetBool("Observe.DaemonSet")
	c.Observe.Deployment = v.GetBool("Observe.Deployment")
	c.Observe.Ingress = v.GetBool("Observe.Ingress")
	c.Observe.Job = v.GetBool("Observe.Job")
	c.Observe.ReplicaSet = v.GetBool("Observe.ReplicaSet")
	c.Observe.ReplicationController = v.GetBool("Observe.ReplicationController")
	c.Observe.Secret = v.GetBool("Observe.AppSecret")
	c.Observe.ServiceAccount = v.GetBool("Observe.ServiceAccount")
	c.Observe.ClusterRole = v.GetBool("Observe.ClusterRole")
	c.Observe.ResourceQuota = v.GetBool("Observe.ResourceQuota")
	c.Observe.NetworkPolicy = v.GetBool("Observe.NetworkPolicy")

	// return configuration as value to avoid thread issues
	// note: does not refresh after config has been loaded though
	return *c, nil
}
