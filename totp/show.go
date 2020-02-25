package totp

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/cacilhas/totp-warehouse/helpers"
	qrcode "github.com/skip2/go-qrcode"
)

func ShowOTP(otp OTP) {
	var file *os.File
	var err error

	if file, err = ioutil.TempFile(os.TempDir(), "*.png"); err != nil {
		helpers.NotifyError(err)
		return
	}
	filename := file.Name()
	defer func() {
		file.Close()
		os.Remove(filename)
	}()

	if err = qrcode.WriteFile(otp.String(), qrcode.High, 512, filename); err != nil {
		helpers.NotifyError(err)
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", filename)
	} else {
		cmd = exec.Command("open", filename)
	}

	if err = cmd.Start(); err != nil {
		helpers.NotifyError(err)
		return
	}

	cmd.Wait()
}
