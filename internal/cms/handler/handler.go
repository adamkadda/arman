// The handler package represents the HTTP layer.
// The objective is to keep handlers nice and narrow. Generally, handlers should
// only interact with (interfaces to) services.
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/adamkadda/arman/pkg/logging"
)

// respondJSON is a convenience function easier struct marshalling and
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
		l.Error("Failed to marshal JSON", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}
