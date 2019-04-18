package pageqr

import (
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/afero"
)

// ImageInfo is a type that expresses info extracted from the parsed image tags
type ImageInfo struct {
	PageURL  string
	ImageSrc string
}

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
			info.PageURL, _ = sel.Attr("title")
			info.ImageSrc, _ = sel.Attr("src")
			run(info)
		})

	return nil
}
