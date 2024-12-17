package server

import (
	"os"
	"os/signal"
	"syscall"
)

// server ...
type server struct {
	log    logger
	stopCh chan os.Signal
	errCh  chan error
}

// Options holds the configuration for the server.
type Options struct {
	Logger logger
}

// Option is a function that configures the server.
type Option func(*server)

// New returns a new server.
func New(options ...Option) *server {
	s := &server{
		stopCh: make(chan os.Signal),
		errCh:  make(chan error),
	}
	for _, option := range options {
		option(s)
	}

	if s.log == nil {
		s.log = NewLogger()
	}

	return s
}

// Start the server.
func (s server) Start() error {
	go func() {
		// Add server startup code here.
		// Send errors to s.errCh.
	}()

	go func() {
		s.stop()
	}()

	s.log.Info("Server started.")
	for {
		select {
		case err := <-s.errCh:
			close(s.errCh)
			return err
		case sig := <-s.stopCh:
			s.log.Info("Server stopped.", "reason", sig.String())
			close(s.stopCh)
			return nil
		}
	}
}

// stop the server.
func (s server) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	// Add server shutdown logic here.

	s.stopCh <- sig
}

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		if options.Logger != nil {
			s.log = options.Logger
		}
	}
}
