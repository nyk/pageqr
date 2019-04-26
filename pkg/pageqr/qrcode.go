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
	"errors"
	"net/url"
	"path"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/afero"
)

// ImageInfo is a type that expresses info extracted from the parsed image tags
type ImageInfo struct {
	PageURL  string
	PagePath string
	ImageSrc string
}

const solidus = 0x2f // solidus is the UTF-8 name for a forward slash

// GenerateImage is a Function to generate a QR code image.
func GenerateImage(siteRoot string, info ImageInfo,
	recovery qrcode.RecoveryLevel, size int) (string, error) {

	fpath, err := FilePathFromURL(siteRoot, info)
	if err != nil {
		return "", err
	}

	file, err := Fs.Create(fpath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	q, err := qrcode.New(info.PageURL, recovery)
	if err != nil {
		return "", err
	}

	err = q.Write(size, file)
	if err != nil {
		return "", err
	}

	return fpath, nil

}

// FilePathFromURL returns a normalized absolute path to write the image file
func FilePathFromURL(siteRoot string, info ImageInfo) (string, error) {

	// Ensure the document root exists
	exists, err := afero.DirExists(Fs, siteRoot)

	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("Document root does not exist")
	}

	urlObject, err := url.Parse(info.ImageSrc)

	if err != nil {
		return "", err
	}

	// Handle absolute URLs
	if urlObject.IsAbs() {
		return filepath.Clean(
			path.Join(
				siteRoot, urlObject.EscapedPath(),
			),
		), nil
	}

	// Handle page relative URLs
	if info.ImageSrc[0] != solidus { // doesn't start with forward slash
		return filepath.Clean(
			path.Join(
				siteRoot, info.PagePath, urlObject.EscapedPath(),
			),
		), nil
	}

	// Site relative urls get handled here
	return filepath.Clean(
		path.Join(
			siteRoot, urlObject.EscapedPath(),
		),
	), nil
}
