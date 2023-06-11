package service

// Option is a function that configures the service.
type Option func(*service)

// Options holds the configuration for the service.
type Options struct {
	Log logger
}

// WithOptions configures the service with the given Options.
func WithOptions(options Options) Option {
	return func(s *service) {
		// Setup on all options from the option struct here.
		s.log = options.Log
	}
}
