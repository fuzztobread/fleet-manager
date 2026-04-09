package vehicle

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler holds the HTTP logic for the vehicle domain.
type Handler struct {
	service *Service
}

func NewHandler(s *Service) http.Handler {
	h := &Handler{service: s}

	r := chi.NewRouter()
	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.Patch("/{id}/status", h.updateStatus)

	return r
}

type createRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

// updateStatusRequest is what we expect in the PATCH body.
type updateStatusRequest struct {
	Status Status `json:"status"`
}

// create handles POST /vehicles
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v, err := h.service.Register(r.Context(), req.Name, req.Capacity)
	if err != nil {
		// service returns plain business errors (e.g. "name is required")
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, v)
}

// list handles GET /vehicles
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch vehicles")
		return
	}

	writeJSON(w, http.StatusOK, vehicles)
}

// get handles GET /vehicles/{id}
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	v, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "vehicle not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not fetch vehicle")
		return
	}

	writeJSON(w, http.StatusOK, v)
}

// updateStatus handles PATCH /vehicles/{id}/status
func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 204 No Content — success, nothing to return
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
