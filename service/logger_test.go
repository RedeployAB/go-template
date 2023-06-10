package service

import (
	"errors"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewDefaultLogger(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		got := NewDefaultLogger()
		want := defaultLogger{out: log.Println}

		if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(defaultLogger{})); diff != "" {
			t.Errorf("NewDefaultLogger() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestDefaultLogger_Info(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		messages := []string{}
		log := defaultLogger{out: testLogFunc(&messages)}

		log.Info("message", "key", "value")
		got := messages[0]

		want := "message=message; key=value"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Info(%s) = unexpected result (-want +got):\n%s\n", want, diff)
		}
	})
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		err := errors.New("error")
		messages := []string{}
		log := defaultLogger{out: testLogFunc(&messages)}

		log.Error(err, "message", "key", "value")
		got := messages[0]

		want := "message=message; error=error; key=value"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Error(%s) = unexpected result (-want +got):\n%s\n", want, diff)
		}
	})
}

func testLogFunc(messages *[]string) func(v ...any) {
	return func(v ...any) {
		*messages = append(*messages, v[0].(string))
	}
}
