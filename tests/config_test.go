package tests

import (
	"regexp"
	"testing"

	"github.com/cacilhas/totp-warehouse/config"
)

func TestConfig(t *testing.T) {
	t.Run("Testing", func(t *testing.T) {
		if !config.Testing() {
			t.Fatalf("failed determinating it's a test")
		}
	})

	t.Run("AppName", func(t *testing.T) {
		expected := "totp-warehouse.test"
		if got := config.AppName(); got != expected {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})

	t.Run("ConfigDir", func(t *testing.T) {
		expected := regexp.MustCompilePOSIX("^/tmp/[0-9]+$")
		if got := config.ConfigDir(); !expected.Match([]byte(got)) {
			t.Fatalf("unexpected config directory %v", got)
		}
	})
}
