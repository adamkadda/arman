package handler

import (
	"errors"
	"net/http"

	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
)

type BiographyHandler struct {
	biographyService *service.BiographyService
}

func NewBiographyHandler(
	biographyService *service.BiographyService,
) *BiographyHandler {
	return &BiographyHandler{
		biographyService: biographyService,
	}
}

// Register registers all biography-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
func (h *BiographyHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /biography/{variant}", h.get)
	mux.HandleFunc("PUT /biography/{variant}", h.update)
}

type biographyRequest struct {
	Content string `json:"content"`
}

func (r *biographyRequest) toDomain(
	variant content.BiographyVariant,
) content.Biography {
	return content.Biography{
		Content: r.Content,
		Variant: variant,
	}
}

type biographyResponse struct {
	Content string `json:"content"`
	Variant string `json:"variant"`
}

func newBiographyResponse(b *content.Biography) biographyResponse {
	return biographyResponse{
		Content: b.Content,
		Variant: string(b.Variant),
	}
}

func (h *BiographyHandler) get(w http.ResponseWriter, r *http.Request) {
	varStr := r.PathValue("variant")

	variant := content.BiographyVariant(varStr)

	biography, err := h.biographyService.Get(r.Context(), variant)
	if err != nil {
		if errors.Is(err, content.ErrInvalidBiographyVariant) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "biography not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newBiographyResponse(biography)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *BiographyHandler) update(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[biographyRequest](w, r)
	if !ok {
		return
	}

	varStr := r.PathValue("variant")

	variant := content.BiographyVariant(varStr)

	biography, err := h.biographyService.Update(r.Context(), req.toDomain(variant))
	if err != nil {
		if errors.Is(err, content.ErrInvalidBiographyVariant) {
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

	resp := newBiographyResponse(biography)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}
