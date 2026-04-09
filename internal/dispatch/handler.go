// internal/dispatch/handler.go
package dispatch

import (
	"encoding/json"
	"errors"
	"net/http"

	"fleet-manager/internal/route"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) http.Handler {
	h := &Handler{svc: svc}

	r := chi.NewRouter()
	r.Post("/jobs", h.handleEnqueue)
	r.Get("/jobs", h.handleQueueStatus)

	return r
}

// POST /dispatch/jobs — enqueue a new job
func (h *Handler) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		From    string  `json:"from"`
		To      string  `json:"to"`
		Urgency Urgency `json:"urgency"`
		MinCap  int     `json:"min_cap"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.svc.Enqueue(body.From, body.To, body.Urgency, body.MinCap)
	if err != nil {
		switch {
		case errors.Is(err, route.ErrNodeNotFound),
			errors.Is(err, route.ErrNoPathExists):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 — accepted, not yet processed
	json.NewEncoder(w).Encode(job)
}

// GET /dispatch/jobs — inspect queue depth and next job
func (h *Handler) handleQueueStatus(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		QueueLen int  `json:"queue_len"`
		Next     *Job `json:"next,omitempty"`
	}{
		QueueLen: h.svc.QueueLen(),
		Next:     h.svc.PeekNext(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
