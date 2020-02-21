package tests

import (
	"regexp"
	"testing"

	"github.com/cacilhas/totp-warehouse/config"
)

func TestConfig(t *testing.T) {
	t.Run("ConfigDir", func(t *testing.T) {
		expected := regexp.MustCompilePOSIX("^/tmp/[0-9]+$")
		if got := config.ConfigDir(); !expected.Match([]byte(got)) {
			t.Fatalf("unexpected config directory %v", got)
		}
	})
}
