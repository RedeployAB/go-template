package server

import (
	"log/slog"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewLogger(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		got := NewLogger()
		want := slog.New(slog.NewJSONHandler(os.Stderr, nil))

		if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(slog.Logger{})); diff != "" {
			t.Errorf("NewLogger() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}
