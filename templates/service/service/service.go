package service

import (
	"os"
	"os/signal"
	"syscall"
)

// service ...
type service struct {
	log    logger
	stopCh chan os.Signal
	errCh  chan error
}

// Options holds the configuration for the service.
type Options struct {
	Logger logger
}

// Option is a function that configures the service.
type Option func(*service)

// New returns a new service.
func New(options ...Option) *service {
	s := &service{
		stopCh: make(chan os.Signal),
		errCh:  make(chan error),
	}
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
	go func() {
		// Add service startup code here.
		// Send errors to s.errCh.
	}()

	go func() {
		s.stop()
	}()

	s.log.Info("Service started.")
	for {
		select {
		case err := <-s.errCh:
			close(s.errCh)
			return err
		case sig := <-s.stopCh:
			s.log.Info("Service stopped.", "reason", sig.String())
			close(s.stopCh)
			return nil
		}
	}
}

// stop the service.
func (s service) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	// Add service shutdown logic here.

	s.stopCh <- sig
}

// WithOptions configures the service with the given Options.
func WithOptions(options Options) Option {
	return func(s *service) {
		if options.Logger != nil {
			s.log = options.Logger
		}
	}
}
