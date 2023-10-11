package main

import (
	"image/jpeg"
	"os"

	"github.com/oarkflow/imaging/imag"
)

func main() {
	f, err := os.Open("./images/sujit-5.jpg")
	if err != nil {
		panic(err)
	}
	d, err := imag.Decode(f, imag.AutoOrientation(true))
	if err != nil {
		panic(err)
	}
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.jpg")
	if err != nil {
		panic(err)
	}
	jpeg.Encode(outputFile, d, nil)
}
