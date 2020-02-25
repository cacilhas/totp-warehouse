package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kirsle/configdir"
)

var (
	appname     string = "totp-warehouse"
	basedir     string
	istest      bool
	appIcon     string
	errorDialog string
	infoDialog  string
	warnDialog  string
)

type Icon int

const (
	// ICON is the application icon path index
	ICON = iota
	// ERROR is the error dialog icon path index
	ERROR
	// INFO is the information dialog icon path index
	INFO
	// WARN is the error warning icon path index
	WARN
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

	appDir := os.Getenv("APPDIR")
	if appDir == "" {
		appIcon = "./assets/key.png"
		errorDialog = "./assets/error.png"
		infoDialog = "./assets/info.png"
		warnDialog = "./assets/warn.png"
	} else {
		appIcon = appDir + "/usr/share/icons/128x128/apps/key.png"
		errorDialog = appDir + "/usr/share/icons/48x48/status/error.png"
		infoDialog = appDir + "/usr/share/icons/48x48/status/info.png"
		warnDialog = appDir + "/usr/share/icons/48x48/status/warn.png"
	}
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

func GetIconPath(index Icon) string {
	switch index {
	case ERROR:
		return errorDialog
	case INFO:
		return infoDialog
	case WARN:
		return warnDialog
	default:
		return appIcon
	}
}
