package event_store

import "context"

type Storage interface {
	Create(ctx context.Context, event EventLog) (string, error)
}
