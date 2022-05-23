package webui

import (
	"encoding/json"
	"net/http"

	"github.com/areknoster/attendgo/domain"
	"github.com/go-chi/chi/v5"
)

var _ Handler = (*AttendeesHandler)(nil)

type JSONAttendee struct {
	ID      string `json:"Id"`
	PhotoID string `json:"PhotoId"`
}

type AttendeesHandler struct {
	storage domain.AtendeeStorage
}

func NewAttendeesHandler(s domain.AtendeeStorage) *AttendeesHandler {
	return &AttendeesHandler{
		storage: s,
	}
}

// Register implements Handler
func (ah *AttendeesHandler) Register(r chi.Router) {
	r.Get("/attendees", func(w http.ResponseWriter, r *http.Request) {
		domainAttendees := ah.storage.List()
		jsonAttendees := make([]JSONAttendee, len(domainAttendees))
		for i, a := range domainAttendees {
			jsonAttendees[i] = JSONAttendee{
				ID:      string(a.ID),
				PhotoID: a.PhotoRef.String(),
			}
		}

		r.Header.Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(jsonAttendees)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
