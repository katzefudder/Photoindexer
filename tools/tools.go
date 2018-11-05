package tools

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const portrait string = "portrait"
const landscape string = "landscape"
const square string = "square"

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

func ImageResize(filename string, outputDir string, suffix string, width int, height int) (int, int, string) {
	filePath, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	imgWidth, imgHeight := GetImageDimension(filename)

	// decode jpeg into image.Image
	img, err := jpeg.Decode(filePath)
	if err != nil {
		log.Fatal(err)
	}
	filePath.Close()

	var newImage image.Image

	var orientation string

	orientation = GetImageOrientation(GetImageDimension(filename))

	switch orientation {
	case "portrait":
		newImage = resize.Resize(uint(width), 0, img, resize.Lanczos3)
	case "landscape":
		newImage = resize.Resize(0, uint(height), img, resize.Lanczos3)
	case "square":
		newImage = resize.Resize(uint(width), 0, img, resize.Lanczos3)
	}

	resultFilePath := filepath.Base(filename)
	resultFileName := strings.TrimSuffix(resultFilePath, filepath.Ext(resultFilePath))

	// create folder if not existing
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	fmt.Println("File: ", resultFilePath, "Width:", imgWidth, "Height:", imgHeight, " Resizing to: ", width, "x", height)

	out, err := os.Create(outputDir + "/" + resultFileName + suffix + filepath.Ext(resultFilePath))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, newImage, nil)
	if err != nil {
		log.Fatal(err)
	}

	return width, height, orientation
}

func GetImageOrientation(imgWidth int, imgHeight int) (orientation string) {
	if imgWidth > imgHeight {
		orientation = "landscape"
	} else if imgWidth < imgHeight {
		orientation = "portrait"
	} else if imgWidth == imgHeight {
		orientation = "square"
	}

	return orientation
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
