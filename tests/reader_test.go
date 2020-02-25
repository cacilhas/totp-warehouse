package tests

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/cacilhas/totp-warehouse/helpers"
	"github.com/cacilhas/totp-warehouse/totp"
)

func TestReader(t *testing.T) {
	t.Run("ReadTotp", func(t *testing.T) {
		t.Run("Hello World", func(t *testing.T) {
			file, _ := os.Open("fixtures/hello.png")
			defer file.Close()
			if got, err := totp.ReadTotp(file); got != "Hello, World!" {
				t.Fatalf("expected Hello, World!, got %v - error: %v", got, err)
			}
		})

		t.Run("not QR code", func(t *testing.T) {
			file, _ := os.Open("fixtures/error.png")
			defer file.Close()
			if got, err := totp.ReadTotp(file); err == nil {
				t.Fatalf("expected error, got %v", got)
			}
		})
	})

	t.Run("ReadTotpFromScreen", func(t *testing.T) {
		t.Run("Hello World", func(t *testing.T) {
			var cmd *exec.Cmd
			var err error
			if cmd, err = helpers.Open("fixtures/hello.png"); err != nil {
				t.Fatalf("could not open fixture: %v", err)
			}
			defer helpers.Kill(cmd)
			sleep, _ := time.ParseDuration("200ms")
			time.Sleep(sleep)
			if got, err := totp.ReadTotpFromScreen(); got != "Hello, World!" {
				t.Fatalf("expected Hello, World!, got %v - error: %v", got, err)
			}
		})
	})
}
