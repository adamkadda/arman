package handler

import (
	"errors"
	"net/http"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
)

// PieceHandler exposes HTTP endpoints for managing pieces.
// It is a thin HTTP-to-service adapter and contains no business logic.
type PieceHandler struct {
	pieceService *service.PieceService
}

func NewPieceHandler(pieceService *service.PieceService) *PieceHandler {
	return &PieceHandler{
		pieceService: pieceService,
	}
}

// Register registers all piece-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
func (h *PieceHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /pieces/{id}", h.get)
	mux.HandleFunc("GET /pieces", h.list)
	mux.HandleFunc("POST /pieces", h.create)
	mux.HandleFunc("PUT /pieces/{id}", h.update)
	mux.HandleFunc("DELETE /pieces/{id}", h.delete)
}

type pieceRequest struct {
	Operation model.Operation `json:"operation"`
	ID        *int            `json:"id"`
	Data      *pieceData      `json:"data"`
}

func (r pieceRequest) Validate() error {
	if err := r.Operation.Validate(); err != nil {
		return err
	}

	if r.Data == nil {
		return model.ErrMissingData
	}

	return nil
}

func (r pieceRequest) toCommand() model.PieceCommand {
	pieceIntent := model.PieceIntent{
		Operation: r.Operation,
		Data:      r.Data.toDomain(r.ID),
	}

	composerIntent := model.ComposerIntent{
		Operation: r.Data.Composer.Operation,
		Data:      r.Data.Composer.Data.toDomain(r.Data.Composer.ID),
	}

	return model.PieceCommand{
		Piece:    pieceIntent,
		Composer: composerIntent,
	}
}

type pieceData struct {
	Title    string          `json:"title"`
	Composer composerRequest `json:"composer"`
}

func (d pieceData) toDomain(id *int) content.Piece {
	piece := content.Piece{
		Title: d.Title,
	}

	if id != nil {
		piece.ID = *id
	}

	return piece
}

type pieceResponse struct {
	ID         int    `json:"piece_id"`
	Title      string `json:"piece_title"`
	ComposerID int    `json:"composer_id"`
}

func newPieceResponse(p *content.Piece) pieceResponse {
	return pieceResponse{
		ID:         p.ID,
		Title:      p.Title,
		ComposerID: p.ComposerID,
	}
}

type pieceWithDetailsResponse struct {
	ID             int    `json:"piece_id"`
	Title          string `json:"piece_title"`
	ComposerID     int    `json:"composer_id"`
	ProgrammeCount int    `json:"programme_count"`
}

func newPieceWithDetailsResponse(
	p *model.PieceWithDetails,
) pieceWithDetailsResponse {
	return pieceWithDetailsResponse{
		ID:             p.Piece.ID,
		Title:          p.Piece.Title,
		ComposerID:     p.Piece.ComposerID,
		ProgrammeCount: p.ProgrammeCount,
	}
}

func (h *PieceHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	piece, err := h.pieceService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "piece not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newPieceResponse(piece)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *PieceHandler) list(w http.ResponseWriter, r *http.Request) {
	pieces, err := h.pieceService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := make([]pieceWithDetailsResponse, len(pieces))
	for i := range pieces {
		resp[i] = newPieceWithDetailsResponse(&pieces[i])
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *PieceHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[pieceRequest](w, r)
	if !ok {
		return
	}

	if err := req.Validate(); err != nil {
		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", err.Error()),
		)
		return
	}

	piece, err := h.pieceService.Create(r.Context(), req.toCommand())
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

	resp := newPieceResponse(piece)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		resp,
	)
}

func (h *PieceHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[pieceRequest](w, r)
	if !ok {
		return
	}

	if err := req.Validate(); err != nil {
		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", err.Error()),
		)
		return
	}

	req.ID = &id

	piece, err := h.pieceService.Update(r.Context(), req.toCommand())
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

	resp := newPieceResponse(piece)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *PieceHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.pieceService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, content.ErrResourceNotFound):
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "piece not found"),
			)
			return
		case errors.Is(err, content.ErrPieceProtected):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "piece in use"),
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
