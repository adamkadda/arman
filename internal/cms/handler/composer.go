package handler

import (
	"errors"
	"net/http"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
)

// ComposerHandler exposes HTTP endpoints for managing composers.
// It is a thin HTTP-to-service adapter and contains no business logic.
type ComposerHandler struct {
	composerService *service.ComposerService
}

func NewComposerHandler(
	composerService *service.ComposerService,
) *ComposerHandler {
	return &ComposerHandler{
		composerService: composerService,
	}
}

// Register registers all composer-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
func (h *ComposerHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /composers/{id}", h.get)
	mux.HandleFunc("GET /composers", h.list)
	mux.HandleFunc("POST /composers", h.create)
	mux.HandleFunc("PUT /composers/{id}", h.update)
	mux.HandleFunc("DELETE /composers/{id}", h.delete)
}

type composerRequest struct {
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
}

func (r *composerRequest) toDomain() content.Composer {
	return content.Composer{
		FullName:  r.FullName,
		ShortName: r.ShortName,
	}
}

func (r *composerRequest) toDomainWithID(id int) content.Composer {
	return content.Composer{
		ID:        id,
		FullName:  r.FullName,
		ShortName: r.ShortName,
	}
}

type composerResponse struct {
	ID        int    `json:"composer_id"`
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
}

func newComposerResponse(c *content.Composer) composerResponse {
	return composerResponse{
		ID:        c.ID,
		FullName:  c.FullName,
		ShortName: c.ShortName,
	}
}

type composerWithDetailsResponse struct {
	ID         int    `json:"composer_id"`
	FullName   string `json:"full_name"`
	ShortName  string `json:"short_name"`
	PieceCount int    `json:"piece_count"`
}

func newComposerWithDetailsResponse(
	c *models.ComposerWithDetails,
) composerWithDetailsResponse {
	return composerWithDetailsResponse{
		ID:         c.Composer.ID,
		FullName:   c.Composer.FullName,
		ShortName:  c.Composer.ShortName,
		PieceCount: c.PieceCount,
	}
}

func (h *ComposerHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	composer, err := h.composerService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "composer not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ComposerHandler) list(w http.ResponseWriter, r *http.Request) {
	composers, err := h.composerService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := make([]composerWithDetailsResponse, len(composers))
	for i := range composers {
		resp[i] = newComposerWithDetailsResponse(&composers[i])
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ComposerHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[composerRequest](w, r)
	if !ok {
		return
	}

	composer, err := h.composerService.Create(r.Context(), req.toDomain())
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

	resp := newComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		resp,
	)
}

func (h *ComposerHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[composerRequest](w, r)
	if !ok {
		return
	}

	composer, err := h.composerService.Update(r.Context(), req.toDomainWithID(id))
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

	resp := newComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ComposerHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.composerService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, content.ErrResourceNotFound):
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "composer not found"),
			)
			return
		case errors.Is(err, content.ErrComposerProtected):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "composer in use"),
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
