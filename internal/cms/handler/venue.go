package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
)

type VenueHandler struct {
	venueService *service.VenueService
}

func (h *VenueHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /venues/{id}", h.get)
	mux.HandleFunc("GET /venues", h.list)
	mux.HandleFunc("POST /venues", h.create)
	mux.HandleFunc("PUT /venues", h.update)
	mux.HandleFunc("DELETE /venues/{id}", h.delete)
}

type VenueResponse struct {
	ID           int    `json:"venue_id"`
	FullAddress  string `json:"full_address"`
	ShortAddress string `json:"short_address"`
}

func NewVenueResponse(v *content.Venue) VenueResponse {
	return VenueResponse{
		ID:           v.ID,
		FullAddress:  v.FullAddress,
		ShortAddress: v.ShortAddress,
	}
}

type VenueWithDetailsResponse struct {
	ID           int    `json:"venue_id"`
	FullAddress  string `json:"full_address"`
	ShortAddress string `json:"short_address"`
	EventCount   int    `json:"event_count"`
}

func NewVenueWithDetailsResponse(
	v *models.VenueWithDetails,
) VenueWithDetailsResponse {
	return VenueWithDetailsResponse{
		ID:           v.Venue.ID,
		FullAddress:  v.Venue.FullAddress,
		ShortAddress: v.Venue.ShortAddress,
		EventCount:   v.EventCount,
	}
}

func (h *VenueHandler) get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.FromContext(r.Context()).Warn(
			"invalid id in path",
			slog.String("id", idStr),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid venue ID",
			},
		)
		return
	}

	venue, err := h.venueService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				map[string]string{
					"error": "venue not found",
				},
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	response := NewVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *VenueHandler) list(w http.ResponseWriter, r *http.Request) {
	venues, err := h.venueService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	response := make([]VenueWithDetailsResponse, len(venues))
	for i, v := range venues {
		response[i] = NewVenueWithDetailsResponse(&v)
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *VenueHandler) create(w http.ResponseWriter, r *http.Request) {
	var v content.Venue
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logging.FromContext(r.Context()).Warn(
			"decode body failed",
			slog.Any("error", err),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid venue in request",
			},
		)
		return
	}

	venue, err := h.venueService.Create(r.Context(), v)
	if err != nil {
		if errors.Is(err, content.ErrInvalidResource) {
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	response := NewVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		response,
	)
}

func (h *VenueHandler) update(w http.ResponseWriter, r *http.Request) {
	var v content.Venue
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logging.FromContext(r.Context()).Warn(
			"decode body failed",
			slog.Any("error", err),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid venue in request",
			},
		)
		return
	}

	venue, err := h.venueService.Update(r.Context(), v)
	if err != nil {
		if errors.Is(err, content.ErrInvalidResource) {
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	response := NewVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *VenueHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.FromContext(r.Context()).Warn(
			"invalid id in path",
			slog.String("id", idStr),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid venue ID",
			},
		)
		return
	}

	err = h.venueService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				map[string]string{
					"error": "venue not found",
				},
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
