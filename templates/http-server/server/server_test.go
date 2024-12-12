package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name  string
		input []Option
		want  *server
	}{
		{
			name:  "default",
			input: []Option{},
			want: &server{
				httpServer: &http.Server{
					Addr:         defaultHost + ":" + defaultPort,
					Handler:      &router{ServeMux: http.NewServeMux()},
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				router: &router{ServeMux: http.NewServeMux()},
				log:    NewLogger(),
			},
		},
		{
			name: "with options",
			input: []Option{
				WithOptions(Options{
					Router:       NewRouter(),
					Logger:       NewLogger(),
					Host:         "localhost",
					Port:         8081,
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				}),
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         "localhost:8081",
					Handler:      &router{ServeMux: http.NewServeMux()},
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				router: &router{ServeMux: http.NewServeMux()},
				log:    NewLogger(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.input...)
			if got == nil {
				t.Errorf("New(%v) = nil; want %v", test.input, test.want)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}), cmpopts.IgnoreUnexported(http.Server{}, http.ServeMux{}, slog.Logger{}), cmpopts.IgnoreFields(server{}, "stopCh", "errCh")); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			log: &mockLogger{
				logs: &logs,
			},
			stopCh: make(chan os.Signal),
			errCh:  make(chan error),
		}
		go func() {
			time.Sleep(time.Millisecond * 100)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		srv.Start()

		want := []string{
			"Server started.",
			"address",
			"localhost:8080",
			"Server stopped.",
			"reason",
			"interrupt",
		}

		if diff := cmp.Diff(want, logs); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestServer_Start_Error(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			log: &mockLogger{
				logs: &logs,
			},
			stopCh: make(chan os.Signal),
			errCh:  make(chan error),
		}

		httpServer := &http.Server{
			Addr: "localhost:8080",
		}

		go func() {
			go func() {
				time.Sleep(time.Millisecond * 100)
				httpServer.Shutdown(context.Background())
			}()
			httpServer.ListenAndServe()
		}()

		time.Sleep(time.Millisecond * 10)
		gotErr := srv.Start()
		if gotErr == nil {
			t.Errorf("Start() = nil; want error")
		}

		wantErr := errors.New("listen tcp 127.0.0.1:8080: bind: address already in use")
		if diff := cmp.Diff(wantErr.Error(), gotErr.Error()); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

type mockLogger struct {
	logs *[]string
}

func (l *mockLogger) Info(msg string, args ...any) {
	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *mockLogger) Error(msg string, args ...any) {
	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}
