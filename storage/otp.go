package storage

import (
	"path"

	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/totp"
	"github.com/cfdrake/go-gdbm"
)

var (
	dbfile string = path.Join(config.ConfigDir(), "storage.db")
)

func init() {
	db, err := gdbm.Open(dbfile, "r")
	if err == nil {
		db.Close()
	} else {
		db, err := gdbm.Open(dbfile, "c")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		db.Sync()
	}
}

func StorageFilename() string {
	return dbfile
}

func SaveOTP(otp totp.OTP) error {
	db, err := gdbm.Open(dbfile, "w")
	if err != nil {
		return err
	}
	defer func() {
		defer db.Close()
		db.Sync()
	}()
	key := otp.Key()
	if db.Exists(key) {
		return db.Replace(key, otp.String())
	}
	return db.Insert(key, otp.String())
}

func RetrieveOTP(key string) (totp.OTP, error) {
	db, err := gdbm.Open(dbfile, "r")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	value, err := db.Fetch(key)
	if err != nil {
		return nil, err
	}
	return totp.Import(value)
}

func DeleteOTP(key string) error {
	db, err := gdbm.Open(dbfile, "w")
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Delete(key)
}

func ListOTPKeys() ([]string, error) {
	db, err := gdbm.Open(dbfile, "r")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	res := make([]string, 0)
	if key, err := db.FirstKey(); err == nil {
		res = append(res, key)
		for {
			key, err = db.NextKey(key)
			if err != nil {
				break
			}
			res = append(res, key)
		}
	}
	return res, nil
}
