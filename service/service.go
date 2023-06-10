package service

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

// service ...
type service struct {
	log logger
}

// New returns a new service.
func New(options ...Option) *service {
	s := &service{}
	for _, option := range options {
		option(s)
	}
	return s.defaults()
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
		s.log.Error(err, "Failed to start service.")
		return err
	case <-time.After(10 * time.Millisecond):
		// Code for when service start is finsihed.
	}

	sig, err := s.shutdown()
	if err != nil {
		s.log.Error(err, "Failed to shutdown service gracefully.")
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

// defaults sets the default values for the service if they are not set.
func (s *service) defaults() *service {
	// Service default setup here.
	if s.log == nil {
		s.log = NewDefaultLogger()
	}
	return s
}
