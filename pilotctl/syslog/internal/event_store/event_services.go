package event_store

import (
	"context"
	"log"
)

var _ Service = &service{}

type service struct {
	storage Storage
}

func NewService(eventStorage Storage) (Service, error) {
	return &service{
		storage: eventStorage,
	}, nil
}

type Service interface {
	Create(ctx context.Context, event EventLog) (string, error)
}

func (s *service) Create(ctx context.Context, event EventLog) (eventID string, err error) {
	log.Println("Save logs to MongoDB")
	eventID, err = s.storage.Create(ctx, event)
	if err != nil {
		return eventID, err
	}
	return eventID, nil
}
