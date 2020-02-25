package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kirsle/configdir"
)

type Icon int

type Config interface {
	AppName() string
	ConfigDir() string
	Testing() bool
	GetIconPath(Icon) string
}

type configT struct {
	appname, basedir, appIcon, errorDialog, infoDialog, warnDialog string
	istest                                                         bool
}

var (
	config = configT{appname: "totp-warehouse"}
)

const (
	// ICON is the application icon path index
	ICON Icon = iota
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
		config.appname = config.appname + ".test"
		config.istest = true
		if config.basedir, err = ioutil.TempDir("/tmp", ""); err != nil {
			panic(err)
		}

	} else {
		config.istest = false
		config.basedir = configdir.LocalConfig(config.appname)
	}
	configdir.MakePath(config.basedir)

	if appDir := os.Getenv("APPDIR"); appDir == "" {
		config.appIcon = "./assets/key.png"
		config.errorDialog = "./assets/error.png"
		config.infoDialog = "./assets/info.png"
		config.warnDialog = "./assets/warn.png"
	} else {
		config.appIcon = appDir + "/usr/share/icons/128x128/apps/key.png"
		config.errorDialog = appDir + "/usr/share/icons/48x48/status/error.png"
		config.infoDialog = appDir + "/usr/share/icons/48x48/status/info.png"
		config.warnDialog = appDir + "/usr/share/icons/48x48/status/warn.png"
	}
}

// GetConfig returns the current configuration instance
func GetConfig() Config {
	return config
}

func (config configT) AppName() string {
	return config.appname
}

func (config configT) ConfigDir() string {
	return config.basedir
}

func (config configT) Testing() bool {
	return config.istest
}

func (config configT) GetIconPath(index Icon) string {
	switch index {
	case ERROR:
		return config.errorDialog
	case INFO:
		return config.infoDialog
	case WARN:
		return config.warnDialog
	default:
		return config.appIcon
	}
}
