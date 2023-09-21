package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewDefaultLogger(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		got := NewDefaultLogger()
		want := slog.New(slog.NewJSONHandler(os.Stderr, nil))

		if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(slog.Logger{})); diff != "" {
			t.Errorf("NewDefaultLogger() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}
