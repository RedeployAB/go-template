package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// server ...
type server struct {
	log logger
}

// Options holds the configuration for the server.
type Options struct {
	Log logger
}

// Option is a function that configures the server.
type Option func(*server)

// New returns a new server.
func New(options ...Option) *server {
	s := &server{}
	for _, option := range options {
		option(s)
	}

	if s.log == nil {
		s.log = NewDefaultLogger()
	}

	return s
}

// Start the server.
func (s server) Start() error {
	errCh := make(chan error, 1)
	go func() {
		// Add server startup code here.
		// Send errors to errCh.
	}()

	select {
	case err := <-errCh:
		s.log.Error("Failed to start server.")
		return err
	case <-time.After(10 * time.Millisecond):
		// Code for when server start is finsihed.
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

	// Add server shutdown logic here.

	return sig, nil
}

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		// Setup on all options from the option struct here.
		s.log = options.Log
	}
}
