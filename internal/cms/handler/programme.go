package handler

import (
	"errors"
	"net/http"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
)

// ProgrammeHandler exposes HTTP endpoints for managing programmes.
// It is a thin HTTP-to-service adapter and contains no business logic.
type ProgrammeHandler struct {
	programmeService *service.ProgrammeService
}

func NewProgrammeHandler(
	programmeService *service.ProgrammeService,
) *ProgrammeHandler {
	return &ProgrammeHandler{
		programmeService: programmeService,
	}
}

// Register registers all programme-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
func (h *ProgrammeHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /programmes/{id}", h.get)
	mux.HandleFunc("GET /programmes", h.list)
	mux.HandleFunc("POST /programmes", h.create)
	mux.HandleFunc("PUT /programmes", h.update)
	mux.HandleFunc("PUT /programmes/{id}/pieces", h.updatePieces)
	mux.HandleFunc("DELETE /programmes/{id}", h.delete)
}

type programmeRequest struct {
	Title string `json:"programme_title"`
}

func (r *programmeRequest) toDomain() content.Programme {
	return content.Programme{
		Title: r.Title,
	}
}

func (r *programmeRequest) toDomainWithID(id int) content.Programme {
	return content.Programme{
		ID:    id,
		Title: r.Title,
	}
}

type programmeResponse struct {
	ID    int    `json:"programme_id"`
	Title string `json:"programme_title"`
}

func newProgrammeResponse(p *content.Programme) programmeResponse {
	return programmeResponse{
		ID:    p.ID,
		Title: p.Title,
	}
}

type programmeWithDetailsResponse struct {
	ID         int    `json:"programme_id"`
	Title      string `json:"programme_title"`
	PieceCount int    `json:"piece_count"`
	EventCount int    `json:"event_count"`
}

func newProgrammeWithDetailsResponse(
	p *model.ProgrammeWithDetails,
) programmeWithDetailsResponse {
	return programmeWithDetailsResponse{
		ID:         p.Programme.ID,
		Title:      p.Programme.Title,
		PieceCount: p.PieceCount,
		EventCount: p.EventCount,
	}
}

type programmePieceResponse struct {
	Title    string `json:"programme_title"`
	Composer string `json:"composer"`
	Sequence int    `json:"sequence"`
}

func newProgrammePieceResponse(pp *content.ProgrammePiece) programmePieceResponse {
	return programmePieceResponse{
		Title:    pp.Piece.Title,
		Composer: pp.Composer.ShortName,
		Sequence: pp.Sequence,
	}
}

type programmeWithPiecesResponse struct {
	ID     int                      `json:"programme_id"`
	Title  string                   `json:"programme_title"`
	Pieces []programmePieceResponse `json:"programmes"`
}

func newProgrammeWithPiecesResponse(
	p *model.ProgrammeWithPieces,
) programmeWithPiecesResponse {
	programmes := make([]programmePieceResponse, len(p.Pieces))
	for i, pp := range p.Pieces {
		programmes[i] = newProgrammePieceResponse(&pp)
	}

	return programmeWithPiecesResponse{
		ID:     p.Programme.ID,
		Title:  p.Programme.Title,
		Pieces: programmes,
	}
}

func (h *ProgrammeHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	programme, err := h.programmeService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "programme not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newProgrammeWithPiecesResponse(programme)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ProgrammeHandler) list(w http.ResponseWriter, r *http.Request) {
	programmes, err := h.programmeService.List(r.Context())
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := make([]programmeWithDetailsResponse, len(programmes))
	for i := range programmes {
		resp[i] = newProgrammeWithDetailsResponse(&programmes[i])
	}

	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ProgrammeHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[programmeRequest](w, r)
	if !ok {
		return
	}

	programme, err := h.programmeService.Create(r.Context(), req.toDomain())
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

	resp := newProgrammeResponse(programme)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		resp,
	)
}

func (h *ProgrammeHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[programmeRequest](w, r)
	if !ok {
		return
	}

	programme, err := h.programmeService.Update(r.Context(), req.toDomainWithID(id))
	if err != nil {
		switch {
		case errors.Is(err, content.ErrInvalidResource):
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", err.Error()),
			)
			return
		case errors.Is(err, content.ErrProgrammeImmutable):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "programme in use"),
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

	resp := newProgrammeResponse(programme)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ProgrammeHandler) updatePieces(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[[]int](w, r)
	if !ok {
		return
	}

	programme, err := h.programmeService.UpdatePieces(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, content.ErrProgrammeImmutable) {
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "programme in use"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newProgrammeWithPiecesResponse(programme)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *ProgrammeHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.programmeService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, content.ErrResourceNotFound):
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "programme not found"),
			)
			return
		case errors.Is(err, content.ErrProgrammeProtected):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "programme in use"),
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
