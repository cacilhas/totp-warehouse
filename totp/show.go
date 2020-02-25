package totp

import (
	"io/ioutil"
	"os"

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

	if cmd, err := helpers.Open(filename); err == nil {
		cmd.Wait()
	} else {
		helpers.NotifyError(err)
	}
}
