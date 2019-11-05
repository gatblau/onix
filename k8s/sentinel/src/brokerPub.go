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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	s "github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

// logs events to kafka
type BrokerPub struct {
	log      *logrus.Entry
	producer s.AsyncProducer
}

func (pub *BrokerPub) Init(c *Config, log *logrus.Entry) {
	pub.log = log

	// creates a broker client configuration
	pub.producer = newProducer(
		strings.Split(c.Publishers.Broker.Brokers, ","),
		&c.Publishers.Broker.Certificate,
		&c.Publishers.Broker.Key,
		&c.Publishers.Broker.CA,
		&c.Publishers.Broker.Verify,
		*log)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range pub.producer.Errors() {
			pub.log.Errorf("Error publishing event: %s", err)
		}
	}()
}

func (pub *BrokerPub) Publish(event Event) {
	// converts the event to json
	bytes, err := toJSON(event)
	if err == nil {
		// creates the producer keyed message so that Kafka sends
		// all messages of the same key to the same partition to guarantee
		// ordered delivery of messages in this case, within a single k8s cluster and namespace
		message := &s.ProducerMessage{
			// message is keyed using the kubernetes cluster:namespace combination
			Key: s.StringEncoder(fmt.Sprintf("%s:%s", event.Change.Host, event.Change.Namespace)),
			// a single topic for all
			Topic: "k8s",
			// the message payload
			Value: s.StringEncoder(string(bytes))}

		// sends the message
		pub.producer.Input() <- message
	} else {
		pub.log.Errorf("Failed to publish event: %s", err)
	}
}

// creates a new async message producer
func newProducer(brokerList []string, certFile *string, keyFile *string, caFile *string, verifySsl *bool, log logrus.Entry) s.AsyncProducer {
	// gets a new configuration with reasonable defaults
	config := s.NewConfig()

	// create a tls configuration if input parameters are provided
	tlsConfig := createTlsConfiguration(certFile, keyFile, caFile, verifySsl)

	// if there is a configuration
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	// the producer will wait for an acknowledgement from the leader only
	// ack 1 -> middle ground between ack 0 and ack all
	config.Producer.RequiredAcks = s.WaitForLocal

	// compresses messages
	config.Producer.Compression = s.CompressionSnappy

	// flushes batches every 500ms
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	// creates the message producer
	producer, err := s.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// log if we're not able to produce messages.
	// note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Errorf("Failed to publish event: %s.", err)
		}
	}()

	return producer
}

// creates a TLS configuration for the message producer
func createTlsConfiguration(certFile *string, keyFile *string, caFile *string, verifySsl *bool) (t *tls.Config) {
	if *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: *verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
