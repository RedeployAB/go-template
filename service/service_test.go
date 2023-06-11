package service

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
		want  *service
	}{
		{
			name:  "default",
			input: []Option{},
			want: &service{
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
			want: &service{
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

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}), cmpopts.IgnoreUnexported(defaultLogger{})); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestService_Start(t *testing.T) {
	t.Run("start service", func(t *testing.T) {
		messages := []string{}
		log := defaultLogger{out: testLogFunc(&messages)}
		srv := &service{
			log: log,
		}
		go func() {
			time.Sleep(time.Millisecond * 100)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		srv.Start()

		want := []string{
			"message=Service shutdown.; reason=interrupt",
		}

		if diff := cmp.Diff(want, messages); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestService_Start_Error(t *testing.T) {

}
