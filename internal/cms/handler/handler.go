// Package handler implements the HTTP layer of the application.
// It is intentionally kept thin and delegates all business logic to
// lower layers.
//
// All routes registered by this package are expected to be protected
// by authentication middleware unless explicitly documented otherwise.
// A small number of performance-related routes are intentionally left
// unauthenticated.
//
// The HTTP layer is not unit-tested in isolation. Handler methods are
// deliberately simple and primarily concerned with request decoding,
// service delegation, and response encoding. Confidence in our handler
// methods is gained through integration-style tests that simulate full
// HTTP request flows against real or test-configured services.
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adamkadda/arman/pkg/logging"
)

// respondJSON is a convenience function for easier struct marshalling and
// preparing responses.
//
// While pre-marshalling adds CPU and memory overhead, it allows us to return
// the appropriate status code in case of an error.
func respondJSON(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(data)
	if err != nil {
		l := logging.FromContext(ctx)
		l.Error("failed to marshal JSON", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}

// pair is a convenience function for adding key-value pairs to respondJSON.
// pair also makes diffs cleaner.
func pair(key, value string) map[string]string {
	return map[string]string{
		key: value,
	}
}

func parseID(
	w http.ResponseWriter,
	r *http.Request,
) (int, bool) {
	idStr := r.PathValue("id")

	if idStr == "" {
		logging.FromContext(r.Context()).Warn("missing id in path")

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", "missing id"),
		)
		return 0, false
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.FromContext(r.Context()).Warn(
			"invalid id in path",
			slog.String("id", idStr),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", "invalid id"),
		)
		return 0, false
	}

	return id, true
}

func parseBody[T any](
	w http.ResponseWriter,
	r *http.Request,
) (T, bool) {
	var req T

	if r.Body == nil {
		logging.FromContext(r.Context()).Warn(
			"request body missing",
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", "missing request body"),
		)
		return req, false
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.FromContext(r.Context()).Warn(
			"decode body failed",
			slog.Any("error", err),
		)

		respondJSON(r.Context(), w,
			http.StatusBadRequest,
			pair("error", "invalid request body"),
		)
		return req, false
	}

	return req, true
}
