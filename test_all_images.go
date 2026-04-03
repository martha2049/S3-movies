package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"

	_ "image/jpeg"
	_ "image/png"
	"golang.org/x/image/webp"
)

func main() {
	
	imagesDir := "images"

	// прочитать все файлы в images
	files, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		fmt.Println("Failed to read images directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(imagesDir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(file.Name(), ": Failed to read file:", err)
			continue
		}

		// для jpeg и png
		_, format, err := image.Decode(bytes.NewReader(data))
		if err == nil {
			fmt.Println(file.Name(), ": Decoded successfully! Format:", format)
			continue
		}

		// для webp
		_, err = webp.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println(file.Name(), ": Failed to decode (jpeg/png/webp):", err)
		} else {
			fmt.Println(file.Name(), ": Decoded successfully as WebP!")
		}
	}
}