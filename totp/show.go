package totp

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/martinlindhe/notify"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	errorDialog string
)

func init() {
	appDir := os.Getenv("APPDIR")
	if appDir == "" {
		errorDialog = "./assets/error.png"
	} else {
		errorDialog = appDir + "/error.png"
	}
}

func ShowOTP(otp OTP) {
	var file *os.File
	var err error

	if file, err = ioutil.TempFile(os.TempDir(), "*.png"); err != nil {
		notifyError(err)
		return
	}
	filename := file.Name()
	defer func() {
		file.Close()
		os.Remove(filename)
	}()

	if err = qrcode.WriteFile(otp.String(), qrcode.High, 512, filename); err != nil {
		notifyError(err)
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", filename)
	} else {
		cmd = exec.Command("open", filename)
	}

	if err = cmd.Start(); err != nil {
		notifyError(err)
		return
	}

	cmd.Wait()
}

func notifyError(err error) {
	notify.Alert("TOTP Warehouse", "Error", err.Error(), errorDialog)
}
