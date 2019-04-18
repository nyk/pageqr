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

				processPage(fpath)
			})
	}
}

func processPage(fpath string) {
	pageqr.ParsePage(fpath, viper.GetString("CssSelector"),
		func(info pageqr.ImageInfo) {
			err := pageqr.GenerateImage(
				info.PageURL, info.ImageSrc,
				viper.Get("recovery").(qrcode.RecoveryLevel), viper.GetInt("ImageSize"))
			if err != nil {
				log.Println(err)
			}
		})
}

func configure() {

	viper.AddConfigPath(".")
	viper.SetConfigName("PageQR")
	viper.SetDefault("CssSelector", "img.qrcode")
	viper.SetDefault("Extensions", []string{".htm", ".html"})
	viper.SetDefault("ImageSize", 256)

	switch viper.GetString("RecoveryLevel") {
	case "Low":
		viper.SetDefault("recovery", qrcode.Low)
	case "Medium":
		viper.SetDefault("recovery", qrcode.Medium)
	case "High":
		viper.SetDefault("recovery", qrcode.High)
	case "Highest":
		viper.SetDefault("recovery", qrcode.Highest)
	default:
		viper.SetDefault("recovery", qrcode.Low)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Configuration error: %s", err)
	}

}
