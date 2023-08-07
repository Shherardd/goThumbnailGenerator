package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
)

func main() {

	//get s.o. type
	//darwin, windows, linux
	osType := runtime.GOOS
	var homeEnv string
	if osType == "darwin" || osType == "linux" {
		homeEnv = "HOME"
	} else if osType == "windows" {
		homeEnv = "USERPROFILE"
	} else {
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}

	sourceDir := filepath.Join(os.Getenv(homeEnv), "Pictures")
	fmt.Println("sourceDir: ", sourceDir)

	dirItems, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	acceptedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	err = os.MkdirAll(filepath.Join(sourceDir, "thumbnails"), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//read all files in sourceDir and filter images
	for _, item := range dirItems {
		if item.IsDir() {
			continue
		}

		//get file extension
		ext := filepath.Ext(item.Name())
		isImage := false
		for _, acceptedExt := range acceptedExtensions {
			if ext == acceptedExt {
				isImage = true
				break
			}
		}

		if isImage {
			tn, err := GenerateThumbnail(sourceDir, item.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = imaging.Save(tn, filepath.Join(sourceDir, "thumbnails", item.Name()))
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Thumbnail generated for: ", item.Name())
		}

	}

}

func GenerateThumbnail(sourceDir string, fileName string) (*image.NRGBA, error) {

	var thumbnail image.Image
	img, err := imaging.Open(filepath.Join(sourceDir, fileName))
	if err != nil {
		return nil, err
	}

	thumbnail = imaging.Thumbnail(img, 100, 100, imaging.NearestNeighbor)

	dst := imaging.New(100, 100, color.NRGBA{0, 0, 0, 0})

	dst = imaging.Paste(dst, thumbnail, image.Pt(0, 0))

	return dst, nil
}
