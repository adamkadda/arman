package handler

import (
	"errors"
	"net/http"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
)

// VenueHandler exposes HTTP endpoints for managing venues.
// It is a thin HTTP-to-service adapter and contains no business logic.
type VenueHandler struct {
	venueService *service.VenueService
}

func NewVenueHandler(
	venueService *service.VenueService,
) *VenueHandler {
	return &VenueHandler{
		venueService: venueService,
	}
}

// Register registers all venue-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
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

func newVenueResponse(v *content.Venue) venueResponse {
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
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	venue, err := h.venueService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "venue not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *VenueHandler) list(w http.ResponseWriter, r *http.Request) {
	venues, err := h.venueService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := make([]venueWithDetailsResponse, len(venues))
	for i := range venues {
		resp[i] = newVenueWithDetailsResponse(&venues[i])
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *VenueHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[venueRequest](w, r)
	if !ok {
		return
	}

	venue, err := h.venueService.Create(r.Context(), req.toDomain())
	if err != nil {
		if errors.Is(err, content.ErrInvalidResource) {
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", err.Error()),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		resp,
	)
}

func (h *VenueHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[venueRequest](w, r)
	if !ok {
		return
	}

	venue, err := h.venueService.Update(r.Context(), req.toDomainWithID(id))
	if err != nil {
		if errors.Is(err, content.ErrInvalidResource) {
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", err.Error()),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newVenueResponse(venue)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *VenueHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.venueService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, content.ErrResourceNotFound):
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "venue not found"),
			)
			return
		case errors.Is(err, content.ErrVenueProtected):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "venue in use"),
			)
			return
		default:
			respondJSON(r.Context(), w,
				http.StatusInternalServerError,
				pair("error", "internal server error"),
			)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
