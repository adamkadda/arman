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

func NewVenueHandler(venueService *service.VenueService) *VenueHandler {
	return &VenueHandler{
		venueService: venueService,
	}
}

func (h *VenueHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /venues/{id}", h.get)
	mux.HandleFunc("GET /venues", h.list)
	mux.HandleFunc("POST /venues", h.create)
	mux.HandleFunc("PUT /venues/{id}", h.update)
	mux.HandleFunc("DELETE /venues/{id}", h.delete)
}

type venueRequest struct {
	Name         string `json:"venue_name"`
	FullAddress  string `json:"full_address"`
	ShortAddress string `json:"short_address"`
}

func (r *venueRequest) toDomain() content.Venue {
	return content.Venue{
		Name:         r.Name,
		FullAddress:  r.FullAddress,
		ShortAddress: r.ShortAddress,
	}
}

func (r *venueRequest) toDomainWithID(id int) content.Venue {
	return content.Venue{
		ID:           id,
		Name:         r.Name,
		FullAddress:  r.FullAddress,
		ShortAddress: r.ShortAddress,
	}
}

type venueResponse struct {
	ID           int    `json:"venue_id"`
	Name         string `json:"venue_name"`
	FullAddress  string `json:"full_address"`
	ShortAddress string `json:"short_address"`
}

func NewVenueResponse(v *content.Venue) venueResponse {
	return venueResponse{
		ID:           v.ID,
		Name:         v.Name,
		FullAddress:  v.FullAddress,
		ShortAddress: v.ShortAddress,
	}
}

type venueWithDetailsResponse struct {
	ID           int    `json:"venue_id"`
	FullAddress  string `json:"full_address"`
	ShortAddress string `json:"short_address"`
	EventCount   int    `json:"event_count"`
}

func newVenueWithDetailsResponse(
	v *models.VenueWithDetails,
) venueWithDetailsResponse {
	return venueWithDetailsResponse{
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

	response := make([]venueWithDetailsResponse, len(venues))
	for i, v := range venues {
		response[i] = newVenueWithDetailsResponse(&v)
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *VenueHandler) create(w http.ResponseWriter, r *http.Request) {
	var req venueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	defer r.Body.Close()

	venue, err := h.venueService.Create(r.Context(), req.toDomain())
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

	var req venueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	defer r.Body.Close()

	venue, err := h.venueService.Update(r.Context(), req.toDomainWithID(id))
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

		if errors.Is(err, content.ErrVenueProtected) {
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				map[string]string{
					"error": "venue referenced by published events",
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
