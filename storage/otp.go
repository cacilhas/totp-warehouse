package storage

import (
	"fmt"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/totp"
)

var (
	dbfile     string        = path.Join(config.ConfigDir(), "data.db")
	dboptions  *bolt.Options = new(bolt.Options)
	bucketname []byte        = []byte("main")
)

func init() {
	timeout, _ := time.ParseDuration("5s")
	dboptions.Timeout = timeout
}

func StorageFilename() string {
	return dbfile
}

func SaveOTP(otp totp.OTP) error {
	var db *bolt.DB
	var err error
	db, err = bolt.Open(dbfile, 0600, dboptions)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		if bucket, err := tx.CreateBucketIfNotExists(bucketname); err == nil {
			return bucket.Put(otp.Key(), otp.Bytes())
		} else {
			return err
		}
	})
}

func RetrieveOTP(key string) (totp.OTP, error) {
	var db *bolt.DB
	var err error
	db, err = bolt.Open(dbfile, 0600, dboptions)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var repr string

	err = db.Update(func(tx *bolt.Tx) error {
		if bucket, err := tx.CreateBucketIfNotExists(bucketname); err == nil {
			res := bucket.Get([]byte(key))
			if res == nil {
				return fmt.Errorf("key %v not found", key)
			}
			repr = string(res)
			return nil

		} else {
			return err
		}
	})

	return totp.Import(repr)
}

func DeleteOTP(key string) error {
	var db *bolt.DB
	var err error
	db, err = bolt.Open(dbfile, 0600, dboptions)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		if bucket, err := tx.CreateBucketIfNotExists(bucketname); err == nil {
			res := bucket.Get([]byte(key))
			if res == nil {
				return fmt.Errorf("key %v not found", key)
			}

			return bucket.Delete([]byte(key))

		} else {
			return err
		}
	})
}

func ListOTPKeys() ([]string, error) {
	var db *bolt.DB
	var err error
	db, err = bolt.Open(dbfile, 0600, dboptions)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var keys []string
	err = db.Update(func(tx *bolt.Tx) error {
		if bucket, err := tx.CreateBucketIfNotExists(bucketname); err == nil {
			res := make([]string, 0)
			err := bucket.ForEach(func(k, _v []byte) error {
				res = append(res, string(k))
				return nil
			})
			keys = res
			return err

		} else {
			return err
		}
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}
