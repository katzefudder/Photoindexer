package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/katzefudder/photoindexer/tools"
)

func main() {
	start := time.Now()

	imageDir := path.Dir(os.Args[1])
	outputDir := path.Dir(os.Args[2])

	tools.IndexFolder(imageDir, outputDir)

	elapsed := time.Since(start)
	log.Printf("indexing took %s for %s files", elapsed, strconv.Itoa(tools.Counter))
}
