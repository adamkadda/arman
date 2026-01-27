// The server package provides an opinionated http server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/adamkadda/arman/pkg/logging"
)

// Server provides a gracefully-stoppable http server implementation. It is safe
// for concurrent use in goroutines.
type Server struct {
	ip       string
	port     string
	listener net.Listener
}

// New creates a new server listening on the provided address that responds to
// the http.Handler. It starts the listener, but does not start the server. If
// an empty port is given, the server randomly chooses one.
func New(port string) (*Server, error) {
	addr := ":" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

// ServeHTTP starts the server and blocks until the provided context is closed.
// When the provided context is closed, the server is gracefully stopped with a
// timeout of 5 seconds.
//
// Once a server has been stopped, it is NOT safe for reuse.
func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	logger := logging.FromContext(ctx)

	errCh := make(chan error, 1)

	// Spawn a goroutine that listens for context closure. When the context is
	// closed, the server is stopped.
	go func() {
		<-ctx.Done()

		logger.Debug("server: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Debug("server: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	// TODO: Create the prometheus metrics proxy.

	// Run the server. This will block until the provided context is closed.
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Debug("server: stopped serving")

	// TODO: Shutdown the prometheus metrics proxy.

	return <-errCh
}

// ServeHTTPHandler is a convenience wrapper that takes an http.Handler.
// It creates a basic http.Server with default settings and calls ServeHTTP.
func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	// TODO: Wrap handler in an OpenCensus handler.

	logging.FromContext(ctx).Info(fmt.Sprintf("listening on port %s...", s.port))

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	}
	return s.ServeHTTP(ctx, srv)
}

// Addr returns the server's listening address in "ip:port" form.
func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

// IP returns the server's IP.
func (s *Server) IP() string {
	return s.ip
}

// Port returns the server's port.
func (s *Server) Port() string {
	return s.port
}
