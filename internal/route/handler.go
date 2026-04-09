package route

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is the HTTP layer for the route domain.
type Handler struct {
	service *Service
}

// NewHandler wires the routes and returns an http.Handler.
func NewHandler(s *Service) http.Handler {
	h := &Handler{service: s}

	r := chi.NewRouter()
	r.Get("/shortest", h.shortest) // GET /routes/shortest?from=KTM&to=JMP
	r.Get("/nodes", h.nodes)       // GET /routes/nodes

	return r
}

// shortest handles GET /routes/shortest?from=KTM&to=JMP
func (h *Handler) shortest(w http.ResponseWriter, r *http.Request) {
	// r.URL.Query().Get() reads ?from= from the URL query string.
	// Returns "" if the param is missing — service will validate that.
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	route, err := h.service.FindRoute(from, to)
	if err != nil {
		// map domain errors to HTTP status codes
		switch {
		case errors.Is(err, ErrNodeNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrNoPathExists):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, route)
}

// nodes handles GET /routes/nodes — returns all valid node IDs.
func (h *Handler) nodes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.Nodes())
}

// helpers — same pattern as vehicle/handler.go
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
