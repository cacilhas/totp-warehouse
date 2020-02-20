package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cacilhas/totp-warehouse/totp"
	"github.com/kirsle/configdir"
)

var (
	basedir    string        // set in init()
	dbfile     string        // set in init()
	dboptions  *bolt.Options = new(bolt.Options)
	bucketname []byte        = []byte("main")
)

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		var err error
		if basedir, err = ioutil.TempDir("/tmp", ""); err != nil {
			panic(err)
		}
	} else {
		basedir = configdir.LocalConfig("totp-warehouse")
	}
	dbfile = path.Join(basedir, "data.db")
	configdir.MakePath(basedir)
	timeout, _ := time.ParseDuration("500ms")
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

// TODO: delete key/OTP

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
