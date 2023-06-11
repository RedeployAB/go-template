package server

import (
	"log"
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
				log: defaultLogger{out: log.Println},
			},
		},
		{
			name: "with options",
			input: []Option{
				WithOptions(Options{
					Log: NewDefaultLogger(),
				}),
			},
			want: &server{
				log: defaultLogger{out: log.Println},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.input...)
			if got == nil {
				t.Errorf("New(%v) = nil; want %v", test.input, test.want)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}), cmpopts.IgnoreUnexported(defaultLogger{})); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		messages := []string{}
		log := defaultLogger{out: testLogFunc(&messages)}
		srv := &server{
			log: log,
		}
		go func() {
			time.Sleep(time.Millisecond * 100)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		srv.Start()

		want := []string{
			"message=Server shutdown.; reason=interrupt",
		}

		if diff := cmp.Diff(want, messages); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestServer_Start_Error(t *testing.T) {

}
