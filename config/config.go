package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kirsle/configdir"
)

var (
	basedir string
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
	configdir.MakePath(basedir)
}

func ConfigDir() string {
	return basedir
}
