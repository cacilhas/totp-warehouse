package helpers

import (
	"github.com/cacilhas/totp-warehouse/config"
	"github.com/martinlindhe/notify"
)

var (
	title = "TOTP Warehouse"
)

func NotifyError(err error) {
	notify.Alert("TOTP Warehouse", "Error", err.Error(), config.GetIconPath(config.ERROR))
}

func NotifyInfo(message string) {
	notify.Notify(title, "Notice", message, config.GetIconPath(config.INFO))
}

func NotifyWarn(message string) {
	notify.Notify(title, "Warning", message, config.GetIconPath(config.WARN))
}
