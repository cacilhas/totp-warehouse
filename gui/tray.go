package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/helpers"
	"github.com/cacilhas/totp-warehouse/storage"
	"github.com/cacilhas/totp-warehouse/totp"
	"github.com/getlantern/systray"
	"github.com/martinlindhe/notify"
)

var (
	icon image.Image
)

func init() {
	var file *os.File
	var err error
	iconPath := config.GetIconPath(config.ICON)
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
	addKeyItem := systray.AddMenuItem("⊕ Add new Key", "")
	systray.AddSeparator()
	fillMenu()
	quitItem := systray.AddMenuItem("⏻ Quit TOTP Warehouse", "")

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
							config.GetIconPath(config.INFO),
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
		for _, key := range keys {
			go dealWithShow(systray.AddMenuItem(fmt.Sprintf("👁 %v", key), "").ClickedCh, key)
			go dealWithGetToken(systray.AddMenuItem(fmt.Sprintf("📄 %v", key), "").ClickedCh, key)
			go dealWithRemove(systray.AddMenuItem(fmt.Sprintf("❌ %v", key), "").ClickedCh, key)
			systray.AddSeparator()
		}
	} else {
		notifyError(err)
	}
}

func notifyError(err error) {
	notify.Alert("TOTP Warehouse", "Error", err.Error(), config.GetIconPath(config.ERROR))
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
						config.GetIconPath(config.INFO),
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
			// TODO: find a dialog lib
			cmd := exec.Command(
				"zenity",
				"--question",
				fmt.Sprintf("--text='Are you sure you want to remove %v?'", key),
			)
			cmd.Start()
			if cmd.Wait() == nil {
				remove(key)
			}
		}
	}
}

func remove(key string) {
	if err := storage.DeleteOTP(key); err == nil {
		notify.Notify(
			"TOTP Warehouse",
			"Notice",
			fmt.Sprintf("%v removed", key),
			config.GetIconPath(config.INFO),
		)
		restart()

	} else {
		notifyError(err)
	}
}

func restart() {
	systray.Quit()
	helpers.Fork(nil)
}

func onExit() {
	//notify.Notify(
	//	"TOTP Warehouse",
	//	"Warn",
	//	"Exiting",
	//	warnDialog,
	//)
}
