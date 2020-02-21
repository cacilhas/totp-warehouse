package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kirsle/configdir"
)

var (
	appname string = "totp-warehouse"
	basedir string
	istest  bool
)

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		var err error
		appname = appname + ".test"
		istest = true
		if basedir, err = ioutil.TempDir("/tmp", ""); err != nil {
			panic(err)
		}

	} else {
		istest = false
		basedir = configdir.LocalConfig(appname)
	}
	configdir.MakePath(basedir)
}

func AppName() string {
	return appname
}

func ConfigDir() string {
	return basedir
}

func Testing() bool {
	return istest
}
