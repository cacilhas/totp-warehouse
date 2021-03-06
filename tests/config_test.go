package tests

import (
	"regexp"
	"testing"

	"github.com/cacilhas/totp-warehouse/config"
)

func TestConfig(t *testing.T) {
	current := config.GetConfig()

	t.Run("Testing", func(t *testing.T) {
		if !current.Testing() {
			t.Fatalf("failed determinating it's a test")
		}
	})

	t.Run("AppName", func(t *testing.T) {
		expected := "totp-warehouse.test"
		if got := current.AppName(); got != expected {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})

	t.Run("ConfigDir", func(t *testing.T) {
		expected := regexp.MustCompilePOSIX("^/tmp/[0-9]+$")
		if got := current.ConfigDir(); !expected.Match([]byte(got)) {
			t.Fatalf("unexpected config directory %v", got)
		}
	})

	t.Run("GetIconPath", func(t *testing.T) {
		t.Run("ICON", func(t *testing.T) {
			expected := "./assets/key.png"
			if got := current.GetIconPath(config.ICON); got != expected {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})

		t.Run("ERROR", func(t *testing.T) {
			expected := "./assets/error.png"
			if got := current.GetIconPath(config.ERROR); got != expected {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})

		t.Run("INFO", func(t *testing.T) {
			expected := "./assets/info.png"
			if got := current.GetIconPath(config.INFO); got != expected {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})

		t.Run("WARN", func(t *testing.T) {
			expected := "./assets/warn.png"
			if got := current.GetIconPath(config.WARN); got != expected {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})
	})
}
