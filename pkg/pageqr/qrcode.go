package pageqr

import (
	qrcode "github.com/skip2/go-qrcode"
)

// Function to generate a QR code image.
func generateImage(content string, imgPath string,
	recovery qrcode.RecoveryLevel, size int) error {
	file, err := Fs.Create(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	q, err := qrcode.New(content, recovery)
	if err != nil {
		return err
	}

	err = q.Write(size, file)
	if err != nil {
		return err
	}

	return nil

}
