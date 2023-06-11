package server

import (
	"net/http"
	"time"
)

// Option is a function that configures the server.
type Option func(*server)

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
