package server

import (
	"net"
	"net/http"
	"strings"
)

// loggingResponseWriter is a wrapper around an http.ResponseWriter that keeps
// track of the status code and length of the response.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// WriteHeader acts as an adapter for the ResponseWriter's WriteHeader method,
// and also keeps track of the status code.
func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write acts as an adapter for the ResponseWriter's Write method,
// and also keeps track of the status code and length of the response.
func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// requestLogger is a middleware that logs the incoming request.
func requestLogger(log logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r)
		log.Info("Request received.", "status", lw.status, "path", r.URL.Path, "method", r.Method, "remoteIp", resolveIP(r))
	})
}

// resolveIP checks request for headers Forwarded, X-Forwarded-For, and X-Real-Ip
// and falls back to the RemoteAddr if none are found.
func resolveIP(r *http.Request) string {
	var addr string
	if f := r.Header.Get("Forwarded"); f != "" {
		for _, segment := range strings.Split(f, ",") {
			addr = strings.TrimPrefix(segment, "for=")
			break
		}
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		addr = strings.Split(xff, ",")[0]
	} else if xrip := r.Header.Get("X-Real-Ip"); xrip != "" {
		addr = xrip
	} else {
		addr = r.RemoteAddr
	}
	ip := strings.Split(addr, ":")[0]
	if net.ParseIP(ip) == nil {
		return "N/A"
	}
	return ip
}
