# Mini Syslog server library
mini syslog server with mongodb as a backend

As standalone mini syslog, config.yaml is required.

Can be used as library as well to integrate with other application

# Standalone use case scenario
1. Compile mini syslog
```shell
#Linux
CGO_ENABLED=0 GOOS=linux go build -o mini-syslog cmd/main/main.go
#MacOS
CGO_ENABLED=0 GOOS=darwin go build -o mini-syslog cmd/main/main.go
```
2. Fill config.yaml file
3. Store mini-syslog executable and config.yaml in same place
4. run ./mini-syslog

><b>Note: </b> Possible to compile that way that config.yaml will be in /etc/mini-syslog/config.yaml and mini-syslog 
> executable will be in /usr/bin/ or /bin/ folder. Persistently run service via systemd.

# Use as a library in existing application code
><b>Info: </b> mini-syslog will use as MongoDB backend to export logs from edge hosts and other services.

- Configuration yaml is parsed to below struct, if config.go is used.
```shell
type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"5514"`
	}
	MongoDB struct {
		Host       string `yaml:"host" env-required:"true"`
		Port       string `yaml:"port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		AuthDB     string `yaml:"auth_db" env-required:"true"`
		Database   string `yaml:"database" env-required:"true"`
		Collection string `yaml:"collection" env-required:"true"`
	} `yaml:"mongodb" env-required:"true"`
}
```
Location: /internal/config.go

- MongoDB client initialization
```shell
import context
host := mngHost
port := mngPort
username := mngUser
password := mngPassword
database := mngDB
mongoClient, err := mongo.NewClient(context.Background(), host, port, username,
		password, database)
if err != nil {
  log.Fatal(err)
}
```
- MongoDB Storage Initialization
```shell
eventStorage := db.NewStorage(mongoClient, cfg.MongoDB.Collection)
eventService, err := event_store.NewService(eventStorage)
if err != nil {
  log.Fatal(err)
}
```
- Mini-syslog Initialization
```shell

BindIP := 0.0.0.0
Type   := port
Port   := 514 #can be any port number - service will listen on tcp and udp ports simultaneously
event := event_store.Event{
	EventService: eventService,
}

listener := event_store.SyslogListener{
	BindIP: cfg.Listen.BindIP,
	Type:   cfg.Listen.Type,
	Port:   cfg.Listen.Port,
}
event.SyslogServer(listener)
```
- Mini Syslog at this moment is listening logs based on RFC3164, which is default in today rsyslog.
><b>Info: </b> Very easily can add other syslog RFC standards. Any logs from anywhere if will be encapsulated 
> in RFC3164 standard and will push to mini-syslog on defined port, mini-syslog will store it in MongoDB collection.
```shell
type RsyslogLogRFC3164 struct {
	Client    string `json:"client" bson:"client"`
	Content   string `json:"content" bson:"content"`
	Facility  int    `json:"facility" bson:"facility"`
	Hostname  string `json:"hostname" bson:"hostname"`
	Priority  int    `json:"priority" bson:"priority"`
	Severity  int    `json:"severity" bson:"severity"`
	Tag       string `json:"tag" bson:"tag"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
	TLSPeer   string `json:"tls_peer" bson:"tls_peer"`
}
```
Location: /internal/event_store/model.go

- MongoDB collection structure, based on this structure logs will be stored in MongoDB
```shell
type EventLog struct {
	EventID         string `json:"event_id,omitempty" yaml:"event_id,omitempty" bson:"event_id,omitempty"`
	Client          string `json:"client,omitempty" yaml:"client,omitempty" bson:"client,omitempty"`
	CreateTimeStamp string `json:"create_time_stamp,omitempty" yaml:"create_time_stamp,omitempty" bson:"create_time_stamp,omitempty"`
	Hostname        string `json:"hostname,omitempty" yaml:"hostname,omitempty" bson:"hostname,omitempty"`
	HostID          string `json:"host_id" yaml:"host_id" bson:"host_id"`
	HostAddress     string `json:"host_address,omitempty" yaml:"host_address,omitempty" bson:"host_address,omitempty"`
	Location        string `json:"location" yaml:"location" bson:"location"`
	Facility        int    `json:"facility,omitempty" yaml:"facility,omitempty" bson:"facility,omitempty"`
	Priority        int    `json:"priority,omitempty" yaml:"priority,omitempty" bson:"priority,omitempty"`
	Severity        int    `json:"severity,omitempty" yaml:"severity,omitempty" bson:"severity,omitempty"`
	Tag             string `json:"tag" yaml:"tag" bson:"tag"`
	EventTimestamp  string `json:"event_timestamp,omitempty" yaml:"event_timestamp,omitempty" bson:"event_timestamp,omitempty"`
	Content         string `json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
	Details         string `json:"details" yaml:"details" bson:"details"`
}
```
Location: /internal/event_store/model.go

><b>Note:</b> Because mini-syslog will run on edge nodes, configure edge nodes OS level rsyslog service which 
> is running as unix socket to send to mini-syslog only critical, emergency, warning and error logs.
```shell
vim /etc/rsyslog.conf
*.=warn;*.=err;*.=emerg;*.=crit  @localhost:514
:wq
systemctl restart rsyslog
```
>Same related to when application is sending logs, restrict to send only critical, emergency, warning and error logs.
> Use in golang for example logrus library and go-syslog to send hook to mini-syslog in proper RFC3164 standard