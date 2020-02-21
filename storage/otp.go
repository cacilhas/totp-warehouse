package storage

import (
	"github.com/99designs/keyring"
	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/totp"
)

var (
	ring keyring.Keyring
)

func init() {
	var err error
	appname := config.AppName()
	kconfig := keyring.Config{
		ServiceName:              appname,
		KeychainName:             "info.cacilhas." + appname,
		KeychainTrustApplication: config.Testing(),
		KeychainSynchronizable:   false,
		FileDir:                  config.ConfigDir(),
		AllowedBackends: []keyring.BackendType{
			keyring.KWalletBackend,
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.PassBackend,
			keyring.FileBackend,
		},
	}
	if config.Testing() {
		kconfig.FilePasswordFunc = func(_p string) (string, error) {
			return "", nil
		}
	}

	if ring, err = keyring.Open(kconfig); err != nil {
		panic(err)
	}
}

func SaveOTP(otp totp.OTP) error {
	return ring.Set(keyring.Item{
		Key:  otp.Key(),
		Data: []byte(otp.String()),
	})
}

func RetrieveOTP(key string) (totp.OTP, error) {
	item, err := ring.Get(key)
	if err != nil {
		return nil, err
	}
	return totp.Import(string(item.Data))
}

func DeleteOTP(key string) error {
	return ring.Remove(key)
}

func ListOTPKeys() ([]string, error) {
	return ring.Keys()
}
