// Copyright © 2019 Nicholas J. Cowham <nyk@cowham.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pageqr

import (
	"crypto/md5"
	"testing"

	"github.com/spf13/afero"

	qrcode "github.com/skip2/go-qrcode"
)

func TestGenerateImage(t *testing.T) {

	Fs.MkdirAll("/var/www/static/images", 0664)

	info := ImageInfo{
		ImageSrc: "/static/images/qrcode.png",
		PagePath: "static/images/test.html",
		PageURL:  "hello world",
	}

	checksum := [16]byte{186, 89, 188, 226, 128, 201, 174, 122, 169, 0, 1, 81, 127, 65, 34, 176}

	fpath, err := GenerateImage("/var/www", info, qrcode.Medium, 256)
	if err != nil {
		t.Error(err)
	}

	fileInfo, err := Fs.Stat(fpath)
	if err != nil {
		t.Error(err)
	}

	if fileInfo.Name() != "qrcode.png" {
		t.Error("File not created")
	}

	file, err := afero.ReadFile(Fs, fpath)
	if err != nil {
		t.Error(err)
	}

	sum := md5.Sum(file)
	if sum != checksum {
		t.Errorf("Checksum: %x does not match expected checksum: %x", sum, checksum)
	}

}

func TestFilePathFromURL(t *testing.T) {

	var siteRoot = "c:/var/www"
	var filePath = "c:\\var\\www\\images\\qrcodes\\test.png"

	Fs.MkdirAll("c:/var/www/images/qrcodes", 0664)

	// Test error for missing document root
	t.Run("siteRootNotFound", func(t *testing.T) {
		_, err := FilePathFromURL("/not/exists", ImageInfo{
			PagePath: "/posts",
			ImageSrc: "/images/qrcodes/test.png",
		})
		if err == nil {
			t.Error("Missing document root error not returned")
		}
	})

	// Test conversion of absolute URLs to file paths
	t.Run("Absolute", func(t *testing.T) {
		fpath, err := FilePathFromURL(siteRoot, ImageInfo{
			PagePath: "/posts",
			ImageSrc: "http://www.example.com/images/qrcodes/test.png",
		})

		if err != nil {
			t.Error(err)
		}

		if fpath != filePath {
			t.Error("Absolute URL failed to convert to file path:", fpath)
		}
	})

	// Test conversion site relative URLs to file paths
	t.Run("SiteRelative", func(t *testing.T) {
		fpath, err := FilePathFromURL(siteRoot, ImageInfo{
			PagePath: "/posts",
			ImageSrc: "/images/qrcodes/test.png",
		})

		if err != nil {
			t.Error(err)
		}

		if fpath != filePath {
			t.Error("Site relative URL failed to convert to file path:", fpath)
		}
	})

	// Test conversion of page relative URLs to file paths
	t.Run("PageRelative", func(t *testing.T) {
		fpath, err := FilePathFromURL(siteRoot, ImageInfo{
			PagePath: "/posts",
			ImageSrc: "images/qrcodes/test.png",
		})

		if err != nil {
			t.Error(err)
		}

		if fpath != "c:\\var\\www\\posts\\images\\qrcodes\\test.png" {
			t.Error("Page relative URL failed to convert to file path:", fpath)
		}
	})

}
