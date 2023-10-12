package main

import (
	"fmt"
	"path/filepath"

	"github.com/oarkflow/imaging/recognizer"
)

func main() {
	var users = map[string]string{
		"amy.jpg":        "Amy",
		"bernadette.jpg": "Bernadette",
		"howard.jpg":     "Howard",
		"penny.jpg":      "Penny",
		"sujit-5.jpg":    "Sujit",
		"raj.jpg":        "Raj",
		"sheldon.jpg":    "Sheldon",
		"leonard.jpg":    "Leonard",
	}
	const fotosDir = "images"

	rec, err := recognizer.New()
	if err != nil {
		panic(err)
	}
	defer rec.Close()
	for img, name := range users {
		err = rec.AddImageToDataset(filepath.Join(fotosDir, img), name)
		if err != nil {
			panic(err)
		}
	}
	rec.SetSamples()
	/*
		images, err := rec.FilterImageById(fotosDir, "Sujit")
		if err != nil {
			panic(err)
		}
		fmt.Println(len(images))
	*/

	images, err := rec.FilterImageByFacesInImage(fotosDir, "elenco3.jpg")
	if err != nil {
		panic(err)
	}
	for _, img := range images {
		fmt.Println(img.Image, img.Faces)
		rec.SaveImage("matching-"+img.Image, img.File)
	}
}
