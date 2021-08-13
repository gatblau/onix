package main

import (
	"context"
	"github.com/gatblau/onix/pilotctl/syslog/internal/config"
	"github.com/gatblau/onix/pilotctl/syslog/internal/event_store"
	"github.com/gatblau/onix/pilotctl/syslog/internal/event_store/db"
	mongo "github.com/gatblau/onix/pilotctl/syslog/pkg/mongodb"
	"log"
)

func main() {
	log.Println("config initializing")
	cfg := config.GetConfig()

	log.Println("mongo client initializing")

	mongoClient, err := mongo.NewClient(context.Background(), cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
		cfg.MongoDB.Password, cfg.MongoDB.Database)
	if err != nil {
		log.Fatal(err)
	}

	eventStorage := db.NewStorage(mongoClient, cfg.MongoDB.Collection)
	eventService, err := event_store.NewService(eventStorage)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("start mini syslog")
	event := event_store.Event{
		EventService: eventService,
	}

	listener := event_store.SyslogListener{
		BindIP: cfg.Listen.BindIP,
		Type:   cfg.Listen.Type,
		Port:   cfg.Listen.Port,
	}
	event.SyslogServer(listener)
}
