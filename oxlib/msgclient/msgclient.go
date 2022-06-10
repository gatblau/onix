package msgclient

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gatblau/onix/artisan/core"
)

//type OnMessageReceived (client mqtt.Client, message mqtt.Message) func
// Msg client client wrapper
type MsgClient struct {
	conf *Conf
	mq   mqtt.Client
}

var (
	mqc *MsgClient
)

type connstatus struct {
	bool
	err error
}

//MsgClient construct a new MsgClient with defaut configuration
func Client() *MsgClient {
	var (
		err    error
		newmqc *MsgClient
	)
	if newmqc == nil {
		newmqc, err = newMsgClient(new(Conf))
		if err != nil {
			log.Fatalf("ERROR: fail to create MsgClient : %s \n", err)
		}
		mqc = newmqc
	}
	return mqc
}

func buildClientOptions(c *Conf) *mqtt.ClientOptions {

	var cid string
	if len(c.getConfOxMsgBrokerClientId()) > 0 {
		cid = c.getConfOxMsgBrokerClientId()
	} else {
		cid = "runner-client"
	}
	clientid := flag.String("clientid", cid, "A clientid for the connection")
	username := flag.String("username", c.getConfOxMsgBrokerUser(), "A username to authenticate to the MQTT server")
	password := flag.String("password", c.getConfOxMsgBrokerPwd(), "Password to match username")
	server := flag.String("server", c.getConfOxMsgBrokerUri(), "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	flag.Parse()
	connOpts := mqtt.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if len(c.getConfOxMsgBrokerUser()) > 0 {
		connOpts.SetUsername(*username)
		if len(c.getConfOxMsgBrokerUser()) > 0 {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: c.getConfOxMsgBrokerInsecureSkipVerify(), ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	return connOpts
}

//func newMsgClient(c *Conf, f OnMessageReceived) (*MsgClient, error) {
func newMsgClient(c *Conf) (*MsgClient, error) {
	connOpts := buildClientOptions(c)
	core.Debug("build client option for mqtt client ")
	client := mqtt.NewClient(connOpts)
	core.Debug("new mqtt client created")
	return &MsgClient{
		conf: c,
		mq:   client,
	}, nil
}

func (mqc *MsgClient) Publish(topic, msg string) error {
	//TODO boolean retain or not to be fixed
	qos := mqc.conf.getConfOxMsgBrokerQoS()
	q := &qos
	if token := mqc.mq.Publish(topic, byte(*q), false, msg); token.Wait() && token.Error() != nil {
		core.Debug("MQ client failed to public message to topic %s : %s\n", topic, token.Error())
		return token.Error()
	} else {
		core.Debug("Published message to topic %s\n", topic)
		return nil
	}
}

func (mqc *MsgClient) Subscribe(handler mqtt.MessageHandler) error {
	qos := mqc.conf.getConfOxMsgBrokerQoS()
	q := &qos
	topic := mqc.conf.getConfOxMsgBrokerTopic()
	if token := mqc.mq.Subscribe(topic, byte(*q), handler); token.Wait() && token.Error() != nil {
		core.Debug("failed to subscribe topic %s : \n %s \n", topic, token.Error())
		return token.Error()
	}
	core.Debug(" subscribed to topic %s \n", topic)
	return nil
}

func (mqc *MsgClient) Start(waitInSeconds int) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// this channel receives a connection
	conn_status := make(chan connstatus, 1)
	// this channel receives a timeout flag
	timeout := make(chan connstatus, 1)

	go func() {
		if token := mqc.mq.Connect(); token.Wait() && token.Error() != nil {
			core.Debug("MQ client failed to connect : %s", token.Error())
			conn_status <- connstatus{bool: false, err: token.Error()}
		} else {
			op := mqc.mq.OptionsReader()
			fmt.Printf("Connected to mqtt broker at [ %s ] with client id [ %s ]\n", mqc.conf.getConfOxMsgBrokerUri(), op.ClientID())
			conn_status <- connstatus{bool: true, err: nil}
		}
	}()

	go func() {
		// timeout period is In Seconds
		time.Sleep(time.Duration(waitInSeconds) * time.Second)
		er := errors.New("mqtt client failed to connect mqtt broker, the timed out period has elapsed\n")
		conn_status <- connstatus{bool: false, err: er}
	}()

	select {
	// the connection has been established before the timeout
	case c := <-conn_status:
		{
			return c.err
		}
	// the connection has not yet returned when the timeout happens
	case t := <-timeout:
		{
			return t.err
		}
	}

	<-stop
	core.Debug("disconnecting mqtt client..")
	mqc.mq.Disconnect(mqc.conf.getConfOxMsgBrokerShutdownGracePeriod())
	core.Debug("mqtt client disconnected..")
	return nil
}

func (mqc *MsgClient) IsConnected() bool {
	return mqc.mq.IsConnected()
}
