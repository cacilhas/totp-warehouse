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
	icon image.Image
)

func init() {
	var iconPath string
	var file *os.File
	var err error
	appDir := os.Getenv("APPDIR")
	if appDir == "" {
		iconPath = "./assets/key.png"
	} else {
		iconPath = appDir + "/key.png"
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
							"notice",
							fmt.Sprintf("%v added", otp),
							"",
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
			label := fmt.Sprintf("Copy %v", key)
			go dealWithGetToken(systray.AddMenuItem(label, "").ClickedCh, key)
		}
		systray.AddSeparator()
	} else {
		notifyError(err)
	}

	if keys, err := storage.ListOTPKeys(); err == nil {
		for _, key := range keys {
			label := fmt.Sprintf("Remove %v", key)
			go dealWithRemove(systray.AddMenuItem(label, "").ClickedCh, key)
		}
		systray.AddSeparator()
	} else {
		notifyError(err)
	}
}

func notifyError(err error) {
	notify.Alert("TOTP Warehouse", "error", err.Error(), "")
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
						"notice",
						fmt.Sprintf("%v copied to clipboard", token),
						"",
					)
				}

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
					"notice",
					fmt.Sprintf("%v removed", key),
					"",
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
	//
}
