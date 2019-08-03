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

	"sync"

	"encoding/json"
	"io/ioutil"

	"github.com/nfnt/resize"
	"github.com/zidizei/iptc"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// Image : a JSON struct
type Image struct {
	Filename         string
	Model            string
	ImageDescription string
	Time             string
}

type Data struct {
	ObjectName               string
	Headline                 string
	Keywords                 []string
	ApplicationRecordVersion int
}

const loglevel string = "info"

const DefaultQuality = 80
const collectMetaData = false

var Counter int

const portrait string = "portrait"
const landscape string = "landscape"
const square string = "square"

func getIptc(imagePath string) *Data {
	info := Data{}

	err := iptc.Load(imagePath, &info)
	if err != nil {
		panic(err)
	}

	return &info
}

// Open : open a file
func Open(filename string) *os.File {
	fp, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	return fp
}

// getExif : return the image's exif
func getExif(filename string) *exif.Exif {
	file := Open(filename)

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	exifData, _ := exif.Decode(file)

	return exifData
}

// writeMetaData : write a file's metadata to JSON
func writeMetaData(filename string, exif *exif.Exif, target string) {

	description, _ := getExifValue(exif, "ImageDescription")
	model, _ := getExifValue(exif, "Model")
	time, _ := getExifValue(exif, "DateTime")

	exifJSON := &Image{
		Filename:         filename,
		Model:            model,
		ImageDescription: description,
		Time:             time,
	}

	file, _ := json.MarshalIndent(exifJSON, "", " ")

	err := ioutil.WriteFile(target, file, 0644)
	if err != nil {
		panic(err)
	}
}

// getExifValue : return the exif's value
func getExifValue(exif *exif.Exif, value exif.FieldName) (tag string, error error) {
	tiffTag, _ := exif.Get(value)

	return fmt.Sprintf("%v", tiffTag), nil
}

// getImageDimension :  get the images' dimension
func getImageDimension(filename string) (int, int) {
	file := Open(filename)
	defer file.Close()

	image, _, err := image.DecodeConfig(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting dimension %s: %v\n", file.Name(), err)
	}
	return image.Width, image.Height
}

// imageResize : resize an image
func imageResize(filename string, outputDir string, suffix string, width int, height int, wg *sync.WaitGroup) (int, int, string) {
	defer wg.Done()
	file := Open(filename)
	defer file.Close()

	imgWidth, imgHeight := getImageDimension(filename)

	if loglevel == "debug" {
		fmt.Println("width: ", imgWidth, "height: ", imgHeight)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while decoding: %s: %v\n", file.Name(), err)
	}

	var newImage image.Image
	var orientation string

	orientation = getImageOrientation(getImageDimension(filename))

	if loglevel == "debug" {
		fmt.Println("orientation: ", orientation)
	}

	switch orientation {
	case "portrait":
		newImage = resize.Resize(uint(width), 0, img, resize.Lanczos3)
	case "landscape":
		newImage = resize.Resize(0, uint(height), img, resize.Lanczos3)
	case "square":
		newImage = resize.Resize(uint(width), 0, img, resize.Lanczos3)
	}

	resultFilePath := filepath.Base(file.Name())
	resultFileName := strings.TrimSuffix(resultFilePath, filepath.Ext(resultFilePath))

	// create folder if not existing
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	if loglevel == "debug" {
		fmt.Println("File: ", resultFilePath, "Width:", imgWidth, "Height:", imgHeight, " Resizing to: ", width, "x", height)
	}

	out, err := os.Create(outputDir + "/" + resultFileName + suffix + filepath.Ext(resultFilePath))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = DefaultQuality

	// write new image to file
	err = jpeg.Encode(out, newImage, &opt)
	if err != nil {
		log.Fatal(err)
	}

	return width, height, orientation
}

// getImageOrientation : get the image's orientation
func getImageOrientation(imgWidth int, imgHeight int) (orientation string) {
	if imgWidth > imgHeight {
		orientation = "landscape"
	} else if imgWidth < imgHeight {
		orientation = "portrait"
	} else if imgWidth == imgHeight {
		orientation = "square"
	}

	return orientation
}

// getFileContentType : get the content type of a file
func getFileContentType(filename string) string {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	file := Open(filename)
	defer file.Close()

	_, err := file.Read(buffer)
	if err != nil {
		panic(err)
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType
}

// IndexFolder : index a folder
func IndexFolder(imageDir string, outputDir string) {

	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		panic(err)
	}

	if loglevel == "debug" {
		fmt.Println("Image dir: ", imageDir)
		fmt.Println("Output dir: ", outputDir)
	}

	for _, file := range files {
		if loglevel == "debug" {
			fmt.Println("Files: ", files)
		}

		actFileName := file.Name()
		actFile := imageDir + "/" + actFileName
		if err != nil {
			panic(err)
		}

		// Get the content
		contentType := getFileContentType(actFile)

		// must be a jpeg, otherwise continue
		if contentType != "image/jpeg" {
			continue
		}

		Counter++

		if collectMetaData == true {
			exif := getExif(actFile)
			writeMetaData(imageDir, exif, outputDir+"/"+actFileName+".json")
		}

		wg := new(sync.WaitGroup)
		wg.Add(3)
		go imageResize(actFile, outputDir, "", 3000, 3000, wg)
		go imageResize(actFile, outputDir+"/med/", "", 1000, 1000, wg)
		go imageResize(actFile, outputDir+"/small/", "", 200, 200, wg)
		wg.Wait()
	}
}
