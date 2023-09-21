package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// service ...
type service struct {
	log logger
}

// Options holds the configuration for the service.
type Options struct {
	Log logger
}

// Option is a function that configures the service.
type Option func(*service)

// New returns a new service.
func New(options ...Option) *service {
	s := &service{}
	for _, option := range options {
		option(s)
	}

	if s.log == nil {
		s.log = NewDefaultLogger()
	}

	return s
}

// Start the service.
func (s service) Start() error {
	errCh := make(chan error, 1)
	go func() {
		// Add service startup code here.
		// Send errors to errCh.
	}()

	select {
	case err := <-errCh:
		s.log.Error("Failed to start service.")
		return err
	case <-time.After(10 * time.Millisecond):
		// Code for when service start is finsihed.
	}

	sig, err := s.shutdown()
	if err != nil {
		s.log.Error("Failed to shutdown service gracefully.")
		return err
	}
	s.log.Info("Service shutdown.", "reason", sig.String())
	return nil
}

// shutdown the service.
func (s service) shutdown() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	// Add service shutdown logic here.

	return sig, nil
}

// WithOptions configures the service with the given Options.
func WithOptions(options Options) Option {
	return func(s *service) {
		// Setup on all options from the option struct here.
		s.log = options.Log
	}
}
