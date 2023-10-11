package main

import (
	"fmt"
	"path/filepath"

	"github.com/oarkflow/imaging/recognizer"
)

func addFile(rec *recognizer.Recognizer, Path, Id string) {

	err := rec.AddImageToDataset(Path, Id)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	const fotosDir = "images"
	const dataDir = "models"

	rec := recognizer.Recognizer{UseGray: false, UseCNN: false, Tolerance: 0.4}
	err := rec.Init(dataDir)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer rec.Close()

	addFile(&rec, filepath.Join(fotosDir, "amy.jpg"), "Amy")
	addFile(&rec, filepath.Join(fotosDir, "bernadette.jpg"), "Bernadette")
	addFile(&rec, filepath.Join(fotosDir, "howard.jpg"), "Howard")
	addFile(&rec, filepath.Join(fotosDir, "penny.jpg"), "Penny")
	addFile(&rec, filepath.Join(fotosDir, "sujit-5.jpg"), "Sujit")
	addFile(&rec, filepath.Join(fotosDir, "raj.jpg"), "Raj")
	addFile(&rec, filepath.Join(fotosDir, "sheldon.jpg"), "Sheldon")
	addFile(&rec, filepath.Join(fotosDir, "leonard.jpg"), "Leonard")

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(filepath.Join(fotosDir, "sujit-1.jpg"))

	if err != nil {
		fmt.Println(err)
		return
	}
	var sujitFace recognizer.Face
	for _, face := range faces {
		if face.Id == "Sujit" {
			sujitFace = face
		}
	}
	img, err := rec.DrawFaces(filepath.Join(fotosDir, "sujit-1.jpg"), []recognizer.Face{sujitFace})

	if err != nil {
		fmt.Println(err)
		return
	}

	rec.SaveImage("faces.jpg", img)

}
