package pageqr

import (
	"crypto/md5"
	"testing"

	"github.com/spf13/afero"

	qrcode "github.com/skip2/go-qrcode"
)

func TestGenerateImage(t *testing.T) {

	Fs.MkdirAll("static/images", 0664)
	testpath := "static/images/qrcode.png"
	checksum := [16]byte{186, 89, 188, 226, 128, 201, 174, 122, 169, 0, 1, 81, 127, 65, 34, 176}

	err := GenerateImage("hello world", testpath, qrcode.Medium, 256)
	if err != nil {
		t.Error(err)
	}

	info, err := Fs.Stat(testpath)
	if err != nil {
		t.Error(err)
	}

	if info.Name() != "qrcode.png" {
		t.Error("File not created")
	}

	file, err := afero.ReadFile(Fs, testpath)
	if err != nil {
		t.Error(err)
	}

	sum := md5.Sum(file)
	if sum != checksum {
		t.Errorf("Checksum: %x does not match expected checksum: %x", sum, checksum)
	}

}
