package tests

import (
	"fmt"
	"testing"

	"github.com/cacilhas/totp-warehouse/totp"
)

func TestGenerator(t *testing.T) {
	var got totp.OTP
	var err error

	t.Run("Import", func(t *testing.T) {
		t.Run("simple URI", func(t *testing.T) {
			uri := "otpauth://totp/Kode%20Code:batalema@cacilhas?secret=ABCDABCD"
			if got, err = totp.Import(uri); err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if issuer := got.Issuer(); issuer != "Kode Code" {
				t.Fatalf("expected Kode Code, got %v", issuer)
			}
			if user := got.User(); user != "batalema@cacilhas" {
				t.Fatalf("expected batalema@cacilhas, got %v", user)
			}
			if length := got.Length(); length != 6 {
				t.Fatalf("expected 6, got %v", length)
			}
			if secret := got.Secret(); secret != "ABCDABCD" {
				t.Fatalf("expected ABCDABCD, got %v", secret)
			}
			if key := got.Key(); key != "batalema@cacilhas@Kode+Code" {
				t.Fatalf("expected batalema@cacilhas@Kode+Code, got %v", key)
			}
			if value := fmt.Sprintf("%v", got); value != uri {
				t.Fatalf("expected %v, got %v", uri, value)
			}
		})

		t.Run("issuer and length supplied", func(t *testing.T) {
			uri := "otpauth://totp/Kode%20Code:batalema@cacilhas?secret=ABCDABCD&issuer=My+Issuer&length=8"
			if got, err = totp.Import(uri); err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if issuer := got.Issuer(); issuer != "My Issuer" {
				t.Fatalf("expected My Issuer, got %v", issuer)
			}
			if user := got.User(); user != "batalema@cacilhas" {
				t.Fatalf("expected batalema@cacilhas, got %v", user)
			}
			if length := got.Length(); length != 8 {
				t.Fatalf("expected 8, got %v", length)
			}
			if key := got.Key(); key != "batalema@cacilhas@My+Issuer" {
				t.Fatalf("expected batalema@cacilhas@My+Issuer, got %v", key)
			}
			if secret := got.Secret(); secret != "ABCDABCD" {
				t.Fatalf("expected ABCDABCD, got %v", secret)
			}
			if value := fmt.Sprintf("%v", got); value != uri {
				t.Fatalf("expected %v, got %v", uri, value)
			}
		})
	})
}
