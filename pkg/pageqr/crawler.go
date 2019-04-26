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
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/afero"
)

// Fs is a package variable to set the Afero filesystem backend
var Fs = afero.NewOsFs()

// CrawlFiles recursively crawls file directories
func CrawlFiles(extensions, excluded StringSet, filepath string, run func(string)) error {

	files, err := afero.ReadDir(Fs, filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fpath := path.Join(filepath, file.Name())

		// skip excluded paths
		if excluded[fpath] == true {
			continue
		}

		// Make a recursive function call on directories
		if file.IsDir() {
			CrawlFiles(extensions, excluded, fpath, run)
			continue
		}

		// Run the callback on files that have the specified file extensions
		if extensions[path.Ext(fpath)] {
			run(fpath)
		}
	}

	return err
}

// ParsePage will parse an HTML page and extract QR Code image tags
func ParsePage(fpath, cssSelect string, run func(ImageInfo)) error {

	f, err := Fs.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return err
	}

	doc.Find(cssSelect).Each(
		func(i int, sel *goquery.Selection) {
			var info ImageInfo
			info.PagePath = fpath
			info.PageURL, _ = sel.Attr("title")
			info.ImageSrc, _ = sel.Attr("src")
			run(info)
		})

	return nil
}
