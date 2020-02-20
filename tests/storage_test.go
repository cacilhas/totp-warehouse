package tests

import (
	"regexp"
	"testing"

	"github.com/cacilhas/totp-warehouse/storage"
	"github.com/cacilhas/totp-warehouse/totp"
)

func TestStorage(t *testing.T) {
	t.Run("storage filename", func(t *testing.T) {
		expected := regexp.MustCompilePOSIX("^/tmp/[0-9]+/data\\.db$")
		if got := storage.StorageFilename(); !expected.Match([]byte(got)) {
			t.Fatalf("unexpected database file %v", got)
		}
	})

	t.Run("list OTP keys", func(t *testing.T) {
		if got, err := storage.ListOTPKeys(); len(got) != 0 || err != nil {
			t.Fatalf("unexpected list or error: %v / %v", got, err)
		}

		totp1, _ := totp.Import("otpauth://totp/Kode%20Code:batalema@cacilhas?secret=ABCDABCD")
		totp2, _ := totp.Import("otpauth://totp/Another:myuser?secret=DEFDEFDEF")
		if err := storage.SaveOTP(totp1); err != nil {
			t.Fatalf("unexpected error storing totp1: %v", err)
		}
		if err := storage.SaveOTP(totp2); err != nil {
			t.Fatalf("unexpected error storing totp2: %v", err)
		}

		if got, err := storage.ListOTPKeys(); len(got) != 2 || err != nil {
			t.Fatalf("unexpected list or error: %v / %v", got, err)
		} else {
			if got[0] != "batalema@cacilhas@Kode+Code" {
				t.Fatalf("expected batalema@cacilhas@Kode+Code, got %v", got[0])
			}
			if got[1] != "myuser@Another" {
				t.Fatalf("expected myuser@Another, got %v", got[0])
			}
		}

		t.Run("retrieve OTP", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				if got, err := storage.RetrieveOTP("myuser@Another"); err == nil {
					if expected := "otpauth://totp/Another:myuser?secret=DEFDEFDEF"; got.String() != expected {
						t.Fatalf("unexpected OTP: %v", got)
					}
				} else {
					t.Fatalf("unexpected error %v", err)
				}
			})

			t.Run("failure", func(t *testing.T) {
				if got, err := storage.RetrieveOTP("unexistent"); err == nil {
					t.Fatalf("should fail: %v", got)
				}
			})
		})
	})
}
