package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// logger is the interface that wraps around methods Info and Error.
type logger interface {
	Info(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
}

// server ...
type server struct {
	log logger
}

// New returns a new server.
func New(options ...Option) *server {
	s := &server{}
	for _, option := range options {
		option(s)
	}
	return s.defaults()
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
		s.log.Error(err, "Failed to start server.")
		return err
	case <-time.After(10 * time.Millisecond):
		// Code for when server start is finsihed.
	}

	sig, err := s.shutdown()
	if err != nil {
		s.log.Error(err, "Failed to shutdown server gracefully.")
		return err
	}
	s.log.Info("Server shutdown.", "reason", sig.String())
	return nil
}

// shutdown the server
func (s server) shutdown() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	// Add server shutdown logic here.

	return sig, nil
}

// defaults sets the default values for the server if they are not set.
func (s *server) defaults() *server {
	// Server default setup here.
	if s.log == nil {
		s.log = NewDefaultLogger()
	}
	return s
}
