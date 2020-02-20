package tests

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

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
			if err := openbrowser("fixtures/hello.png"); err != nil {
				t.Fatalf("could not open fixture: %v", err)
			}
			time.Sleep(time.Second * 1)
			if got, err := totp.ReadTotpFromScreen(); got != "Hello, World!" {
				t.Fatalf("expected Hello, World!, got %v - error: %v", got, err)
			}
		})
	})
}

func openbrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
