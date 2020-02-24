package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cacilhas/totp-warehouse/storage"
	"github.com/cacilhas/totp-warehouse/totp"
	"github.com/getlantern/systray"
	"github.com/martinlindhe/notify"
)

var (
	icon        image.Image
	errorDialog string
	infoDialog  string
	warnDialog  string
)

func init() {
	var iconPath string
	var file *os.File
	var err error
	appDir := os.Getenv("APPDIR")
	if appDir == "" {
		iconPath = "./assets/key.png"
		errorDialog = "./assets/error.png"
		infoDialog = "./assets/info.png"
		warnDialog = "./assets/warn.png"
	} else {
		iconPath = appDir + "/usr/share/icons/128x128/apps/key.png"
		errorDialog = appDir + "/usr/share/icons/48x48/status/error.png"
		infoDialog = appDir + "/usr/share/icons/48x48/status/info.png"
		warnDialog = appDir + "/usr/share/icons/48x48/status/warn.png"
	}
	if file, err = os.Open(iconPath); err != nil {
		panic(err)
	}
	defer file.Close()
	if icon, err = png.Decode(file); err != nil {
		panic(err)
	}
}

func Start() {
	systray.Run(onReady, onExit)
}

func onReady() {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, icon); err != nil {
		panic(err)
	}
	systray.SetIcon(buf.Bytes())
	systray.SetTitle("TOTP Warehouse")
	systray.SetTooltip("TOTP Warehouse")
	addKeyItem := systray.AddMenuItem("Add new Key", "")
	systray.AddSeparator()
	fillMenu()
	quitItem := systray.AddMenuItem("Quit TOTP Warehouse", "")

	for {
		select {
		case <-addKeyItem.ClickedCh:
			if code, err := totp.ReadTotpFromScreen(); err == nil {
				if otp, err := totp.Import(code); err == nil {
					if err := storage.SaveOTP(otp); err == nil {
						notify.Notify(
							"TOTP Warehouse",
							"Notice",
							fmt.Sprintf("%v added", otp),
							infoDialog,
						)
						restart()

					} else {
						notifyError(err)
					}

				} else {
					notifyError(err)
				}

			} else {
				notifyError(err)
			}
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func fillMenu() {
	if keys, err := storage.ListOTPKeys(); err == nil {
		if len(keys) > 0 {
			for _, key := range keys {
				label := fmt.Sprintf("Show %v", key)
				go dealWithShow(systray.AddMenuItem(label, "").ClickedCh, key)
			}
			systray.AddSeparator()
			for _, key := range keys {
				label := fmt.Sprintf("Copy %v", key)
				go dealWithGetToken(systray.AddMenuItem(label, "").ClickedCh, key)
			}
			systray.AddSeparator()
			for _, key := range keys {
				label := fmt.Sprintf("Remove %v", key)
				go dealWithRemove(systray.AddMenuItem(label, "").ClickedCh, key)
			}
			systray.AddSeparator()
		}
	} else {
		notifyError(err)
	}
}

func notifyError(err error) {
	notify.Alert("TOTP Warehouse", "Error", err.Error(), errorDialog)
}

func dealWithGetToken(channel <-chan struct{}, key string) {
	for {
		select {
		case <-channel:
			if otp, err := storage.RetrieveOTP(key); err == nil {
				token := otp.Token()
				if err := clipboard.WriteAll(token); err == nil {
					notify.Notify(
						"TOTP Warehouse",
						"Notice",
						fmt.Sprintf("%v copied to clipboard", token),
						infoDialog,
					)
				}

			} else {
				notifyError(err)
			}
		}
	}
}

func dealWithShow(channel <-chan struct{}, key string) {
	for {
		select {
		case <-channel:
			if otp, err := storage.RetrieveOTP(key); err == nil {
				totp.ShowOTP(otp)

			} else {
				notifyError(err)
			}
		}
	}
}

func dealWithRemove(channel <-chan struct{}, key string) {
	for {
		select {
		case <-channel:
			if err := storage.DeleteOTP(key); err == nil {
				notify.Notify(
					"TOTP Warehouse",
					"Notice",
					fmt.Sprintf("%v removed", key),
					infoDialog,
				)
				restart()

			} else {
				notifyError(err)
			}
		}
	}
}

func restart() {
	systray.Quit()
	exec.Command(os.Args[0], os.Args[1:]...).Start()
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func onExit() {
	notify.Notify(
		"TOTP Warehouse",
		"Warn",
		"Exiting",
		warnDialog,
	)
}
