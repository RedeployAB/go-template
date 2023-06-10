package server

// Option is a function that configures the server.
type Option func(*server)

// Options holds the configuration for the server.
type Options struct {
	Log logger
}

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		// Setup on all options from the option struct here.
		s.log = options.Log
	}
}
