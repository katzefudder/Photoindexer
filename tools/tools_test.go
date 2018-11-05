package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileDimensions(t *testing.T) {
	width, height := GetImageDimension("../images/landscape.jpg")
	assert.Equal(t, width, 4000, "width should be 4000")
	assert.Equal(t, height, 2667, "width should be 2667")
}

func M(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}

func TestGetFileContentTypeFileFails(t *testing.T) {
	var a, b int
	a, b = GetImageDimension("../images/gibtsnet")
	fmt.Println(a, b)
	//assert.Panics(t, , "The code did not panic")
}

func TestPortraitImageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, err := filepath.Abs("../images/portrait.jpg")
	if err != nil {
		panic(err)
	}

	width, height, orientation = ImageResize(absPath, "./testing", "", 100, 100)
	assert.Equal(t, width, 100, "width should be 100")
	assert.Equal(t, height, 100, "height should be 100")
	assert.Equal(t, orientation, "portrait", "orientation should be portrait")
}

func TestSquareImageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, err := filepath.Abs("../images/square.jpg")
	if err != nil {
		panic(err)
	}

	width, height, orientation = ImageResize(absPath, "./testing", "", 100, 100)
	assert.Equal(t, width, 100, "width should be 100")
	assert.Equal(t, height, 100, "height should be 100")
	assert.Equal(t, orientation, "square", "orientation should be square")
}

func TestLandscapeImageResize(t *testing.T) {
	var width, height int
	var orientation string

	absPath, err := filepath.Abs("../images/landscape.jpg")
	if err != nil {
		panic(err)
	}

	width, height, orientation = ImageResize(absPath, "./testing", "", 3000, 3000)
	assert.Equal(t, width, 3000, "width should be 3000")
	assert.Equal(t, height, 3000, "height should be 2000")
	assert.Equal(t, orientation, "landscape", "orientation should be landscape")
}

func TestGetFileContentType(t *testing.T) {
	file, err := os.Open("../images/portrait.jpg")
	contentType, err := GetFileContentType(file)
	assert.Equal(t, contentType, "image/jpeg", "should be image/jpeg")
	assert.Equal(t, err, nil, "should be no error")
}

func TestGetFileContentTypeError(t *testing.T) {
	file, err := os.Open("../images/gibtsnet.jpg")
	contentType, err := GetFileContentType(file)
	assert.Equal(t, contentType, "", "should be ''")
	assert.Equal(t, err.Error(), "invalid argument", "should be 'invalid argument'")
}
