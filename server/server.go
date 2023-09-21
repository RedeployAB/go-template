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

// server holds an http.Server, a router and it's configured options.
type server struct {
	httpServer *http.Server
	router     *http.ServeMux
	log        logger
}

// Options holds the configuration for the server.
type Options struct {
	Router       *http.ServeMux
	Log          logger
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Option is a function that configures the server.
type Option func(*server)

// New returns a new server.
func New(options ...Option) *server {
	s := &server{
		httpServer: &http.Server{
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
	}
	for _, option := range options {
		option(s)
	}

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

	return s
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
		s.log.Error("Failed to start server.")
		return err
	case <-time.After(10 * time.Millisecond):
		s.log.Info("Server started.", "address", s.httpServer.Addr)
	}

	sig, err := s.shutdown()
	if err != nil {
		s.log.Error("Failed to shutdown server gracefully.")
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

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		s.router = options.Router
		s.log = options.Log

		s.httpServer.Handler = options.Router
		s.httpServer.Addr = options.Host + ":" + options.Port
		s.httpServer.ReadTimeout = options.ReadTimeout
		s.httpServer.WriteTimeout = options.WriteTimeout
		s.httpServer.IdleTimeout = options.IdleTimeout
	}
}
