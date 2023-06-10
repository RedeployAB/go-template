package server

import (
	"context"
	"errors"
	"log"
	"net/http"
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
					Handler:      &http.ServeMux{},
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				router: &http.ServeMux{},
				log:    defaultLogger{out: log.Println},
			},
		},
		{
			name: "with options",
			input: []Option{
				WithOptions(Options{
					Router:       http.NewServeMux(),
					Log:          NewDefaultLogger(),
					Host:         "localhost",
					Port:         "8081",
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				}),
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         "localhost:8081",
					Handler:      &http.ServeMux{},
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				router: &http.ServeMux{},
				log:    defaultLogger{out: log.Println},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.input...)
			if got == nil {
				t.Errorf("New(%v) = nil; want %v", test.input, test.want)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}), cmpopts.IgnoreUnexported(http.Server{}, http.ServeMux{}, defaultLogger{})); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		logger := mockLogger{
			logs: &logs,
		}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			log: logger,
		}
		go func() {
			time.Sleep(time.Millisecond * 100)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		srv.Start()

		want := []string{
			"Server started.;address;localhost:8080",
			"Server shutdown.;reason;interrupt",
		}

		if diff := cmp.Diff(want, logs); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestServer_Start_Error(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		logger := mockLogger{
			logs: &logs,
		}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			log: logger,
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
