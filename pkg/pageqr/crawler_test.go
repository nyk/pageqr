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
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
)

const numFiles = 9

var extensions StringSet

func TestMain(m *testing.M) {

	Fs = afero.NewMemMapFs()

	Fs.MkdirAll("/static/html", 0664)
	Fs.MkdirAll("/static/html/posts/jan", 0664)
	Fs.Mkdir("/static/html/posts/feb", 0664)

	for i := 1; i <= numFiles; i++ {
		s := fmt.Sprintf(
			`<!DOCTYPE html><html><head><title>This is test page %d</title></head>
				<body><h1>Test Page %d</h1>
					<img src='/images/qrcodes/qr%d.png' title='test%d.html' class='qrcode'>
					<p>That's all folks</p>
				<body/>
			</html>`, i, i, i, i)
		afero.WriteFile(Fs, fmt.Sprintf("/static/html/test%d.html", i), []byte(s), 0664)
		afero.WriteFile(Fs, fmt.Sprintf("/static/html/posts/test%d.html", i), []byte(s), 0664)
		afero.WriteFile(Fs, fmt.Sprintf("/static/html/posts/jan/test%d.html", i), []byte(s), 0664)
		afero.WriteFile(Fs, fmt.Sprintf("/static/html/posts/feb/test%d.html", i), []byte(s), 0664)
	}

	extensions = NewStringSet([]string{".html", ".htm"})

	os.Exit(m.Run())

}

func TestParsePage(t *testing.T) {

	for i := 1; i <= numFiles; i++ {
		testFile := fmt.Sprintf("/static/html/test%d.html", i)
		err := ParsePage(testFile, "img.qrcode",
			func(info ImageInfo) {
				if info.PagePath != fmt.Sprintf("/static/html/test%d.html", i) {
					t.Error("Parsing error: ", info.PagePath)
				}

				if info.ImageSrc != fmt.Sprintf("/images/qrcodes/qr%d.png", i) {
					t.Error("Parsing error: ", info.ImageSrc)
				}
			})

		if err != nil {
			t.Error(err)
		}
	}
}

func TestParsePageMissingFileError(t *testing.T) {

	err := ParsePage("/static/html/none.html", "img.qrcode",
		func(info ImageInfo) {})

	if err == nil {
		t.Error("Expected file not found error not thrown!")
	}

}

func TestCrawlFiles(t *testing.T) {

	excluded := NewStringSet([]string{"/static/html/posts/jan"})
	var counter = 1

	CrawlFiles(extensions, excluded, "/static/html/posts/feb",
		func(fpath string) {
			if fpath != fmt.Sprintf("/static/html/posts/feb/test%d.html", counter) {
				t.Error(fpath)
			}
			counter++
		})

}

func TestCrawlRecurse(t *testing.T) {
	var counter = 0
	CrawlFiles(extensions, StringSet{}, "/static/html/posts",
		func(fpath string) {
			counter++
		})

	if counter != 3*numFiles {
		t.Error("Crawler recursion test failed")
	}
}
