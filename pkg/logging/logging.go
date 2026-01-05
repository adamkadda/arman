// The logging package initializes and configures logging.
package logging

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// NewLogger creates a new logger with the given configuration.
func NewLogger(level string, development bool) *slog.Logger {
	var handler slog.Handler

	if development {
		options := slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler = slog.NewTextHandler(os.Stderr, &options)
	} else {
		options := slog.HandlerOptions{
			Level: levelToSlogLevel(level),
		}
		handler = slog.NewJSONHandler(os.Stderr, &options)
	}

	return slog.New(handler)
}

// NewLoggerFromEnv creates a new logger from the environment. It consumes LOG_LEVEL
// for determining the level, and LOG_MODE for determining the handler options.
func NewLoggerFromEnv() *slog.Logger {
	level := os.Getenv("LOG_LEVEL")
	development := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "development"
	return NewLogger(level, development)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger exists,
// a default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

// levelToSlogLevel converts the given string to the corresponding slog level value.
func levelToSlogLevel(s string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarn:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}

type loggingWriter struct {
	http.ResponseWriter
	logger     *slog.Logger
	statusCode int
	size       int
}

// WriteHeader is a wrapper around the net/http implementation, but it also
// keeps track of the first status code written.
func (w *loggingWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write is a wrapper around the net/http implementation, but it also ensures that
// we always have a valid status code in our logginWriter's statusCode field.
// As Write attempts to write the response body, it will catch and log any errors.
//
// Lastly, Write keeps track of how many bytes it (attempted to) write. Since this
// method can be called multiple times, we use the increment operator to update
// the loggingWriter's size field.
func (w *loggingWriter) Write(b []byte) (int, error) {

	// NOTE: Edge case: WriteHeader not called. i.e. implicit 200 OK
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	if err != nil {
		w.logger.Error(
			"Failed to write response body",
			slog.String("err", err.Error()),
		)
	}

	w.size += n

	return n, err
}

// The logging package Middleware uses the default logger initialized once by the
// DefaultLogger function. Middleware assigns a requestID to each request it
// intercepts, and logs a Debug-level log for the request's start and completion.
//
// Additionally, it wraps the http.ResponseWriter around the loggingWriter to keep
// track of how many bytes were written and the status code returned.
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := uuid.NewString()
			w.Header().Set("X-Request-ID", requestID)

			logger := DefaultLogger().With(
				slog.Group("request",
					slog.String("id", requestID),
					slog.String("method", r.Method),
					slog.String("uri", r.URL.RequestURI()),
					slog.String("ip", r.RemoteAddr),
				),
			)

			ctx := WithLogger(r.Context(), logger)
			r = r.WithContext(ctx)

			logger.Debug("Request started")

			lw := &loggingWriter{
				ResponseWriter: w,
				logger:         logger,
			}

			next.ServeHTTP(lw, r)

			logger.Debug("Request completed",
				slog.Group("response",
					slog.Int("status", lw.statusCode),
					slog.Int("size", lw.size),
				),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
