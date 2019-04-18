package main

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/nyk/PageQR/pkg/pageqr"
	qrcode "github.com/skip2/go-qrcode"
)

func main() {

	configure()

	// Create a lookup map of file paths to exclude.
	excluded := pageqr.NewStringSet(viper.GetStringSlice("Exclude"))
	extensions := pageqr.NewStringSet(viper.GetStringSlice("Extensions"))

	// Start harvesting files and processing them.
	for _, dir := range viper.GetStringSlice("Include") {
		pageqr.CrawlFiles(extensions, excluded, dir,
			func(fpath string) {
				fpath, err := filepath.Abs(fpath)
				if err != nil {
					return
				}
				log.Println("Yay: ", fpath)
				pageqr.ParsePage(fpath, viper.GetString("CssSelector"),
					func(info pageqr.ImageInfo) {
						log.Println(info.PageURL)
					})
			})
	}
}

func configure() {

	viper.AddConfigPath(".")
	viper.SetConfigName("PageQR")
	viper.SetDefault("CssSelector", "img.qrcode")
	viper.SetDefault("Extensions", []string{".htm", ".html"})
	viper.SetDefault("RecoveryLevel", qrcode.Medium)
	viper.SetDefault("PixelSize", 256)

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Configuration error: %s", err)
	}

}
