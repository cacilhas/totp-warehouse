package totp

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/liyue201/goqr"
	"github.com/vova616/screenshot"
)

func ReadTotpFromScreen() (string, error) {
	image, err := screenshot.CaptureScreen()
	if err != nil {
		return "", err
	}

	var file *os.File

	if file, err = ioutil.TempFile(os.TempDir(), "*.png"); err != nil {
		return "could not open file", err
	}
	filename := file.Name()
	defer func() {
		file.Close()
		os.Remove(filename)
	}()

	if err = png.Encode(file, image); err != nil {
		return "could not encode PNG", err
	}

	file.Close()
	if file, err = os.Open(filename); err != nil {
		return "error trying to reopen temp file", err
	}
	return ReadTotp(file)
}

func ReadTotp(file *os.File) (string, error) {
	var img image.Image
	var err error

	if img, _, err = image.Decode(file); err != nil {
		return "error decoding image", err
	}

	qrcodes, err := goqr.Recognize(img)
	if err != nil {
		return "error decoding QR code", err
	}

	for _, matrix := range qrcodes {
		return string(matrix.Payload), nil
	}

	return "", fmt.Errorf("no QR code found")
}
