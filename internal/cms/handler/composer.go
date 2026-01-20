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

type ComposerHandler struct {
	composerService *service.ComposerService
}

func NewComposerHandler(composerService *service.ComposerService) *ComposerHandler {
	return &ComposerHandler{
		composerService: composerService,
	}
}

func (h *ComposerHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /composers/{id}", h.get)
	mux.HandleFunc("GET /composers", h.list)
	mux.HandleFunc("POST /composers", h.create)
	mux.HandleFunc("PUT /composers", h.update)
	mux.HandleFunc("DELETE /composers/{id}", h.delete)
}

type ComposerResponse struct {
	ID        int    `json:"composer_id"`
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
}

func NewComposerResponse(c *content.Composer) ComposerResponse {
	return ComposerResponse{
		ID:        c.ID,
		FullName:  c.FullName,
		ShortName: c.ShortName,
	}
}

type ComposerWithDetailsResponse struct {
	ID         int    `json:"composer_id"`
	FullName   string `json:"full_name"`
	ShortName  string `json:"short_name"`
	PieceCount int    `json:"piece_count"`
}

func NewComposerWithDetailsResponse(
	c *models.ComposerWithDetails,
) ComposerWithDetailsResponse {
	return ComposerWithDetailsResponse{
		ID:         c.Composer.ID,
		FullName:   c.Composer.FullName,
		ShortName:  c.Composer.ShortName,
		PieceCount: c.PieceCount,
	}
}

func (h *ComposerHandler) get(w http.ResponseWriter, r *http.Request) {
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
				"error": "invalid composer ID",
			},
		)
		return
	}

	composer, err := h.composerService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				map[string]string{
					"error": "composer not found",
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

	response := NewComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *ComposerHandler) list(w http.ResponseWriter, r *http.Request) {
	composers, err := h.composerService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	response := make([]ComposerWithDetailsResponse, len(composers))
	for i, v := range composers {
		response[i] = NewComposerWithDetailsResponse(&v)
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *ComposerHandler) create(w http.ResponseWriter, r *http.Request) {
	var v content.Composer
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logging.FromContext(r.Context()).Warn(
			"decode body failed",
			slog.Any("error", err),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid composer in request",
			},
		)
		return
	}

	composer, err := h.composerService.Create(r.Context(), v)
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

	response := NewComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		response,
	)
}

func (h *ComposerHandler) update(w http.ResponseWriter, r *http.Request) {
	var v content.Composer
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logging.FromContext(r.Context()).Warn(
			"decode body failed",
			slog.Any("error", err),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid composer in request",
			},
		)
		return
	}

	composer, err := h.composerService.Update(r.Context(), v)
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

	response := NewComposerResponse(composer)
	respondJSON(r.Context(), w,
		http.StatusOK,
		response,
	)
}

func (h *ComposerHandler) delete(w http.ResponseWriter, r *http.Request) {
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
				"error": "invalid composer ID",
			},
		)
		return
	}

	err = h.composerService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				map[string]string{
					"error": "composer not found",
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
