package server

import (
	"log/slog"
	"os"
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
				log: NewLogger(),
			},
		},
		{
			name: "with options",
			input: []Option{
				WithOptions(Options{
					Logger: NewLogger(),
				}),
			},
			want: &server{
				log: NewLogger(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.input...)
			if got == nil {
				t.Errorf("New(%v) = nil; want %v", test.input, test.want)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}), cmpopts.IgnoreUnexported(slog.Logger{}), cmpopts.IgnoreFields(server{}, "stopCh", "errCh")); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		srv := &server{
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
			"Server stopped.",
			"reason",
			"interrupt",
		}

		if diff := cmp.Diff(want, logs); diff != "" {
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
		messages = append(messages, v.(string))
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *mockLogger) Error(msg string, args ...any) {
	messages := []string{msg}
	for _, v := range args {
		messages = append(messages, v.(string))
	}
	*l.logs = append(*l.logs, messages...)
}
