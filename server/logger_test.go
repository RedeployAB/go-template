package server

import (
	"log"
	"strings"
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
	msg := ""
	msgPtr := &msg
	l := defaultLogger{out: testLogFunc(msgPtr)}

	l.Info("message", "key", "value")
	got := msg

	want := "message=message; key=value"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Info(%s) = unexpected result (-want +got):\n%s\n", want, diff)
	}

}

func testLogFunc(msg *string) func(v ...any) {
	return func(v ...any) {
		*msg = v[0].(string)
	}
}

type mockLogger struct {
	logs *[]string
}

func (l mockLogger) Info(msg string, keysAndValues ...any) {
	m := make([]string, len(keysAndValues)+1)
	m[0] = msg
	for i := 1; i < len(m); i++ {
		m[i] = keysAndValues[i-1].(string)
	}
	*l.logs = append(*l.logs, strings.Join(m, ";"))
}

func (l mockLogger) Error(err error, msg string, keysAndValues ...any) {
	m := make([]string, len(keysAndValues)+2)
	m[0] = err.Error()
	m[1] = msg
	for i := 2; i < len(m); i++ {
		m[i] = keysAndValues[i-2].(string)
	}
	*l.logs = append(*l.logs, strings.Join(m, ";"))
}
