package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	tools "photoindexer/tools"
)

var BaseDir = ""

func main() {

	BaseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	imageDir := path.Dir(os.Args[1])
	outputDir := path.Dir(os.Args[2])

	files, err := ioutil.ReadDir(imageDir)
	fmt.Println("Processing dir: ", imageDir)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		file, err := os.Open(imageDir + "/" + file.Name())

		if err != nil {
			panic(err)
		}

		// Get the content
		contentType, err := tools.GetFileContentType(file)
		if err != nil {
			panic(err)
		}

		// must be a jpeg
		if contentType != "image/jpeg" {
			continue
		}

		absPath, err := filepath.Abs(file.Name())
		if err != nil {
			panic(err)
		}

		// metadata.GetIptc(absPath)
		// metadata.GetExif(absPath)

		tools.ImageResize(absPath, BaseDir+"/"+outputDir, "", 3000, 3000)
		tools.ImageResize(absPath, BaseDir+"/"+outputDir+"/med/", "", 1000, 1000)
		tools.ImageResize(absPath, BaseDir+"/"+outputDir+"/small/", "", 200, 200)
	}
}
