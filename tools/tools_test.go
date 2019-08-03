package tools

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileDimensions(t *testing.T) {
	filename, _ := filepath.Abs("../images/landscape.jpg")
	width, height := getImageDimension(filename)
	assert.Equal(t, width, 4000, "width should be 4000")
	assert.Equal(t, height, 2667, "width should be 2667")
}

func TestPortraitImageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, _ := filepath.Abs("../images/portrait.jpg")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	width, height, orientation = imageResize(absPath, "./testing", "", 100, 100, wg)
	assert.Equal(t, width, 100, "width should be 100")
	assert.Equal(t, height, 100, "height should be 100")
	assert.Equal(t, orientation, "portrait", "orientation should be portrait")
}

func TestSquareImageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, _ := filepath.Abs("../images/square.jpg")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	width, height, orientation = imageResize(absPath, "./testing", "", 100, 100, wg)
	assert.Equal(t, width, 100, "width should be 100")
	assert.Equal(t, height, 100, "height should be 100")
	assert.Equal(t, orientation, "square", "orientation should be square")
}

func TestLandscapeimageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, _ := filepath.Abs("../images/landscape.jpg")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	width, height, orientation = imageResize(absPath, "./testing", "", 3000, 3000, wg)
	assert.Equal(t, width, 3000, "width should be 3000")
	assert.Equal(t, height, 3000, "height should be 2000")
	assert.Equal(t, orientation, "landscape", "orientation should be landscape")
}

func TestGetFileContentType(t *testing.T) {
	absPath, _ := filepath.Abs("../images/portrait.jpg")
	contentType := getFileContentType(absPath)
	assert.Equal(t, contentType, "image/jpeg", "should be image/jpeg")
}

func TestGetImageOrientation(t *testing.T) {
	absPath, _ := filepath.Abs("../images/portrait.jpg")
	imgWidth, imgHeight := getImageDimension(absPath)
	orientation := getImageOrientation(imgWidth, imgHeight)
	assert.Equal(t, orientation, "portrait")

	absPath, _ = filepath.Abs("../images/landscape.jpg")
	imgWidth, imgHeight = getImageDimension(absPath)
	orientation = getImageOrientation(imgWidth, imgHeight)
	assert.Equal(t, orientation, "landscape")
}

func TestOpenNotExistingFile(t *testing.T) {
	absPath, _ := filepath.Abs("../images/gibtsnet.jpg")
	assert.Panics(t, func() { Open(absPath) })
}

func TestGetIptcFromJpg(t *testing.T) {
	absPath, _ := filepath.Abs("../images/landscape.jpg")
	data := getIptc(absPath)
	fmt.Printf("%#v\n", data)
	assert.Equal(t, data.ObjectName, "DEL2 Playoff-Viertelfinale - Spiel 4 - Rote Teufel EC Bad Nauhei")
}

func TestGetIPTCFromNotExistingFile(t *testing.T) {
	absPath, _ := filepath.Abs("../images/gibtsnet.jpg")
	assert.Panics(t, func() { getIptc(absPath) })
}

func TestGetIptc(t *testing.T) {
	// TODO: impelement
}
