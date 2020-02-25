package helpers

import (
	"github.com/cacilhas/totp-warehouse/config"
	"github.com/martinlindhe/notify"
)

var (
	title         = "TOTP Warehouse"
	currentConfig = config.GetConfig()
)

func NotifyError(err error) {
	notify.Alert(title, "Error", err.Error(), currentConfig.GetIconPath(config.ERROR))
}

func NotifyInfo(message string) {
	notify.Notify(title, "Notice", message, currentConfig.GetIconPath(config.INFO))
}

func NotifyWarn(message string) {
	notify.Notify(title, "Warning", message, currentConfig.GetIconPath(config.WARN))
}
