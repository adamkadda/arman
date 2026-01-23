package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/service"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
)

// EventHandler exposes HTTP endpoints for managing events.
// It is a thin HTTP-to-service adapter and contains no business logic.
type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(
	eventService *service.EventService,
) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// Register registers all event-related HTTP routes on the provided ServeMux.
// Routes are registered at the root and assume JSON request and response bodies.
func (h *EventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /events/{id}", h.get)
	mux.HandleFunc("GET /events", h.list)
	mux.HandleFunc("POST /events", h.create)
	mux.HandleFunc("PUT /events/{id}", h.update)
	mux.HandleFunc("PUT /events/{id}/notes", h.updatesNotes)
	mux.HandleFunc("PUT /events/{id}/draft", h.draft)
	mux.HandleFunc("PUT /events/{id}/publish", h.publish)
	mux.HandleFunc("PUT /events/{id}/archive", h.archive)
	mux.HandleFunc("DELETE /events/{id}", h.delete)
}

type eventRequest struct {
	Title       string     `json:"title"`
	Date        *time.Time `json:"date"`
	TicketLink  *string    `json:"ticket_link"`
	VenueID     *int       `json:"venue_id"`
	ProgrammeID *int       `json:"programme_id"`
}

func (r *eventRequest) toDomain() content.Event {
	return content.Event{
		Title:       r.Title,
		Date:        r.Date,
		TicketLink:  r.TicketLink,
		VenueID:     r.VenueID,
		ProgrammeID: r.ProgrammeID,
	}
}

func (r *eventRequest) toDomainWithID(id int) content.Event {
	return content.Event{
		ID:          id,
		Title:       r.Title,
		Date:        r.Date,
		TicketLink:  r.TicketLink,
		VenueID:     r.VenueID,
		ProgrammeID: r.ProgrammeID,
	}
}

type eventResponse struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Date        *time.Time     `json:"date"`
	TicketLink  *string        `json:"ticket_link"`
	VenueID     *int           `json:"venue_id"`
	ProgrammeID *int           `json:"programme_id"`
	Status      content.Status `json:"status"`
	Notes       *string        `json:"notes"`
}

func newEventResponse(e *content.Event) eventResponse {
	return eventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Date:        e.Date,
		TicketLink:  e.TicketLink,
		VenueID:     e.VenueID,
		ProgrammeID: e.ProgrammeID,
		Status:      e.Status,
		Notes:       e.Notes,
	}
}

type eventWithTimestampsResponse struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Date        *time.Time     `json:"date"`
	TicketLink  *string        `json:"ticket_link"`
	VenueID     *int           `json:"venue_id"`
	ProgrammeID *int           `json:"programme_id"`
	Status      content.Status `json:"status"`
	Notes       *string        `json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func newEventWithTimestampsResponse(
	e *models.EventWithTimestamps,
) eventWithTimestampsResponse {
	return eventWithTimestampsResponse{
		ID:          e.Event.ID,
		Title:       e.Event.Title,
		Date:        e.Event.Date,
		TicketLink:  e.Event.TicketLink,
		VenueID:     e.Event.VenueID,
		ProgrammeID: e.Event.ProgrammeID,
		Status:      e.Event.Status,
		Notes:       e.Event.Notes,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

type eventWithProgrammeResponse struct {
	ID          int                          `json:"id"`
	Title       string                       `json:"title"`
	Date        *time.Time                   `json:"date"`
	TicketLink  *string                      `json:"ticket_link"`
	VenueID     *int                         `json:"venue_id"`
	ProgrammeID *int                         `json:"programme_id"`
	Status      content.Status               `json:"status"`
	Notes       *string                      `json:"notes"`
	Programme   *programmeWithPiecesResponse `json:"programme"`
}

func newEventWithProgrammeResponse(
	e *models.EventWithProgramme,
) eventWithProgrammeResponse {
	programme := newProgrammeWithPiecesResponse(e.Programme)
	return eventWithProgrammeResponse{
		ID:          e.Event.ID,
		Title:       e.Event.Title,
		Date:        e.Event.Date,
		TicketLink:  e.Event.TicketLink,
		VenueID:     e.Event.VenueID,
		ProgrammeID: e.Event.ProgrammeID,
		Status:      e.Event.Status,
		Notes:       e.Event.Notes,
		Programme:   &programme,
	}
}

type notesRequest struct {
	Notes string `json:"notes"`
}

func (h *EventHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	event, err := h.eventService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrResourceNotFound) {
			respondJSON(r.Context(), w,
				http.StatusNotFound,
				pair("error", "event not found"),
			)
			return
		}

		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newEventWithProgrammeResponse(event)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *EventHandler) list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var status *content.Status
	val, ok := query["status"]
	if ok && len(val) > 0 && val[0] != "" {
		s := content.Status(val[0])
		status = &s
	}

	var timeframe *content.Timeframe
	val, ok = query["timeframe"]
	if ok && len(val) > 0 && val[0] != "" {
		s := content.Timeframe(val[0])
		timeframe = &s
	}

	detailed := false
	val, ok = query["detailed"]
	if ok && len(val) > 0 && val[0] != "" {
		b, err := strconv.ParseBool(val[0])
		if err != nil {
			logging.FromContext(r.Context()).Warn(
				"invalid 'detailed' parameter",
				slog.String("detailed", val[0]),
			)

			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", "invalid 'detailed' parameter"),
			)
			return
		}

		detailed = b
	}

	var events any
	var err error

	ctx := r.Context()

	if detailed {
		events, err = h.eventService.ListWithTimestamp(ctx, status, timeframe)
	} else {
		events, err = h.eventService.List(ctx, status, timeframe)
	}

	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	buildResponse := func(events any) any {
		switch e := events.(type) {
		case []content.Event:
			resp := make([]eventResponse, len(e))
			for i := range e {
				resp[i] = newEventResponse(&e[i])
			}
			return resp
		case []models.EventWithTimestamps:
			resp := make([]eventWithTimestampsResponse, len(e))
			for i := range e {
				resp[i] = newEventWithTimestampsResponse(&e[i])
			}
			return resp
		}
		return nil
	}

	respondJSON(ctx, w,
		http.StatusOK,
		buildResponse(events),
	)
}

func (h *EventHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := parseBody[eventRequest](w, r)
	if !ok {
		return
	}

	event, err := h.eventService.Create(r.Context(), req.toDomain())
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

	resp := newEventResponse(event)
	respondJSON(r.Context(), w,
		http.StatusCreated,
		resp,
	)
}

func (h *EventHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[eventRequest](w, r)
	if !ok {
		return
	}

	event, err := h.eventService.Update(r.Context(), req.toDomainWithID(id))
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

	resp := newEventWithProgrammeResponse(event)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *EventHandler) updatesNotes(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	req, ok := parseBody[notesRequest](w, r)
	if !ok {
		return
	}

	event, err := h.eventService.UpdateNotes(r.Context(), id, req.Notes)
	if err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	resp := newEventResponse(event)
	respondJSON(r.Context(), w,
		http.StatusOK,
		resp,
	)
}

func (h *EventHandler) draft(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.eventService.Draft(r.Context(), id); err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) publish(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.eventService.Publish(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, content.ErrInvalidResource):
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", err.Error()),
			)
			return
		case errors.Is(err, content.ErrProgrammeHasNoPieces):
			respondJSON(r.Context(), w,
				http.StatusBadRequest,
				pair("error", err.Error()),
			)
			return
		case errors.Is(err, content.ErrEventNotPublishable):
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", err.Error()),
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

func (h *EventHandler) archive(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.eventService.Archive(r.Context(), id); err != nil {
		respondJSON(r.Context(), w,
			http.StatusInternalServerError,
			pair("error", "internal server error"),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.eventService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, content.ErrEventProtected) {
			respondJSON(r.Context(), w,
				http.StatusForbidden,
				pair("error", "published event protected"))
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
