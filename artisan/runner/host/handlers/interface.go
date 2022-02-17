package handlers

import (
	"net/http"
)

type EventHandler interface {
    HandleEvent(w http.ResponseWriter, r *http.Request)
}

 