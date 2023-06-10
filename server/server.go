package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Defaults for server configuration.
const (
	defaultHost         = "0.0.0.0"
	defaultPort         = "8080"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

// logger is the interface that wraps around methods Info and Error.
type logger interface {
	Info(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
}

// server holds an http.Server, a router and it's configured options.
type server struct {
	httpServer *http.Server
	router     *http.ServeMux
	log        logger
}

// New returns a new server.
func New(options ...Option) *server {
	s := &server{
		httpServer: &http.Server{},
	}
	for _, option := range options {
		option(s)
	}
	return s.defaults()
}

// Start the server.
func (s server) Start() error {
	s.routes()

	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
	}()

	select {
	case err := <-errCh:
		s.log.Error(err, "Failed to start server.")
		return err
	case <-time.After(10 * time.Millisecond):
		s.log.Info("Server started.", "address", s.httpServer.Addr)
	}

	sig, err := s.shutdown()
	if err != nil {
		s.log.Error(err, "Failed to shutdown server gracefully.")
		return err
	}
	s.log.Info("Server shutdown.", "reason", sig.String())
	return nil
}

// shutdown the server.
func (s server) shutdown() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return nil, err
	}
	return sig, nil
}

// defaults sets the default values for the server if they are not set.
func (s *server) defaults() *server {
	if s.router == nil {
		s.router = http.NewServeMux()
		s.httpServer.Handler = s.router
	}
	if s.log == nil {
		s.log = NewDefaultLogger()
	}
	if len(s.httpServer.Addr) == 0 {
		s.httpServer.Addr = defaultHost + ":" + defaultPort
	}
	if s.httpServer.ReadTimeout == 0 {
		s.httpServer.ReadTimeout = defaultReadTimeout
	}
	if s.httpServer.WriteTimeout == 0 {
		s.httpServer.WriteTimeout = defaultWriteTimeout
	}
	if s.httpServer.IdleTimeout == 0 {
		s.httpServer.IdleTimeout = defaultIdleTimeout
	}

	return s
}
