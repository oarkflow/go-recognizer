package main

import (
	"os"

	"github.com/oarkflow/imaging/recognizer"
)

func main() {
	f, err := os.Open("./images/sujit-1.jpg")
	if err != nil {
		panic(err)
	}
	recognizer.AlignImage(f)
}
