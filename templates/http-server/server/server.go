package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	router     *router
	tls        TLSConfig
	log        logger
	stopCh     chan os.Signal
	errCh      chan error
}

// TLSConfig holds the configuration for the server's TLS settings.
type TLSConfig struct {
	Certificate string
	Key         string
}

// isEmpty returns true if the TLSConfig is empty.
func (c TLSConfig) isEmpty() bool {
	return len(c.Certificate) == 0 && len(c.Key) == 0
}

// Options holds the configuration for the server.
type Options struct {
	Router       *router
	TLSConfig    TLSConfig
	Logger       logger
	Host         string
	Port         int
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
		stopCh: make(chan os.Signal),
		errCh:  make(chan error),
	}
	for _, option := range options {
		option(s)
	}

	if s.router == nil {
		s.router = NewRouter()
	}
	if s.log == nil {
		s.log = NewLogger()
	}
	if len(s.httpServer.Addr) == 0 {
		s.httpServer.Addr = defaultHost + ":" + defaultPort
	}
	if s.httpServer.Handler == nil {
		s.httpServer.Handler = s.router
	}

	return s
}

// Start the server.
func (s server) Start() error {
	s.routes()

	go func() {
		if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errCh <- err
		}
	}()

	go func() {
		s.stop()
	}()

	s.log.Info("Server started.", "address", s.httpServer.Addr)
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

// listenAndServe wraps around http.Server ListenAndServe and
// ListenAndServeTLS depending on TLS configuration.
func (s *server) listenAndServe() error {
	if !s.tls.isEmpty() {
		s.httpServer.TLSConfig = newTLSConfig()
		return s.httpServer.ListenAndServeTLS(s.tls.Certificate, s.tls.Key)
	}
	return s.httpServer.ListenAndServe()
}

// stop the server.
func (s server) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.errCh <- err
	}

	s.stopCh <- sig
}

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		if options.Router != nil {
			s.router = options.Router
			s.httpServer.Handler = s.router
		}
		if !options.TLSConfig.isEmpty() {
			s.tls = options.TLSConfig
		}
		if options.Logger != nil {
			s.log = options.Logger
		}
		if len(options.Host) > 0 || options.Port > 0 {
			s.httpServer.Addr = options.Host + ":" + strconv.Itoa(options.Port)
		}
		if options.ReadTimeout > 0 {
			s.httpServer.ReadTimeout = options.ReadTimeout
		}
		if options.WriteTimeout > 0 {
			s.httpServer.WriteTimeout = options.WriteTimeout
		}
		if options.IdleTimeout > 0 {
			s.httpServer.IdleTimeout = options.IdleTimeout
		}
	}
}

// newTLSConfig returns a new tls.Config.
func newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
