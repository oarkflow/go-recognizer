package recognizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"

	goFace "github.com/oarkflow/imaging/go-face"
)

// Data descriptor of the human face.
type Data struct {
	Id         string
	Descriptor goFace.Descriptor
}

// Face holds coordinates and descriptor of the human face.
type Face struct {
	Data
	Rectangle image.Rectangle
}

type Option struct {
	Tolerance float32
	UseCNN    bool
	UseGray   bool
	ModelDir  string
}

/*
A Recognizer creates face descriptors for provided images and
classifies them into categories.
*/
type Recognizer struct {
	opt     *Option
	rec     *goFace.Recognizer
	dataset []Data
}

func New(opt ...*Option) (*Recognizer, error) {
	cfg := &Option{UseGray: false, UseCNN: false}
	if len(opt) > 0 {
		cfg = opt[0]
	}
	if cfg.Tolerance == 0 {
		cfg.Tolerance = 0.4
	}
	if cfg.ModelDir == "" {
		cfg.ModelDir = "models"
	}
	rec := &Recognizer{
		opt:     cfg,
		dataset: make([]Data, 0),
	}
	r, err := goFace.NewRecognizer(cfg.ModelDir)
	if err == nil {
		rec.rec = r
	}
	return rec, err
}

/*
Close frees resources taken by the Recognizer. Safe to call multiple
times. Don't use Recognizer after close call.
*/
func (_this *Recognizer) Close() {
	_this.rec.Close()
}

/*
AddImageToDataset add a sample image to the dataset
*/
func (_this *Recognizer) AddImageToDataset(Path string, Id string) error {
	file := Path
	var err error
	if _this.opt.UseGray {
		file, err = _this.createTempGrayFile(file, Id)
		if err != nil {
			return err
		}
		defer os.Remove(file)
	}
	var faces []goFace.Face
	if _this.opt.UseCNN {
		faces, err = _this.rec.RecognizeFileCNN(file)
	} else {
		faces, err = _this.rec.RecognizeFile(file)
	}
	if err != nil {
		return err
	}
	if len(faces) == 0 {
		return errors.New("Not a face on the image")
	}
	if len(faces) > 1 {
		return errors.New("Not a single face on the image")
	}
	f := Data{}
	f.Id = Id
	f.Descriptor = faces[0].Descriptor
	_this.dataset = append(_this.dataset, f)
	return nil
}

/*
SetSamples sets known descriptors so you can classify the new ones.
*/
func (_this *Recognizer) SetSamples() {
	var samples []goFace.Descriptor
	var avengers []int32
	for i, f := range _this.dataset {
		samples = append(samples, f.Descriptor)
		avengers = append(avengers, int32(i))
	}
	_this.rec.SetSamples(samples, avengers)
}

/*
RecognizeSingle returns face if it's the only face on the image or nil otherwise.
Only JPEG format is currently supported.
*/
func (_this *Recognizer) RecognizeSingle(Path string) (goFace.Face, error) {
	file := Path
	var err error
	if _this.opt.UseGray {
		file, err = _this.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")
		if err != nil {
			return goFace.Face{}, err
		}
		defer os.Remove(file)
	}
	var idFace *goFace.Face
	if _this.opt.UseCNN {
		idFace, err = _this.rec.RecognizeSingleFileCNN(file)
	} else {
		idFace, err = _this.rec.RecognizeSingleFile(file)
	}
	if err != nil {
		return goFace.Face{}, fmt.Errorf("Can't recognize: %v", err)

	}
	if idFace == nil {
		return goFace.Face{}, fmt.Errorf("Not a single face on the image")
	}
	return *idFace, nil
}

/*
RecognizeMultiples returns all faces found on the provided image, sorted from
left to right. Empty list is returned if there are no faces, error is
returned if there was some error while decoding/processing image.
Only JPEG format is currently supported.
*/
func (_this *Recognizer) RecognizeMultiples(Path string) ([]goFace.Face, error) {
	file := Path
	var err error
	if _this.opt.UseGray {
		file, err = _this.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")
		if err != nil {
			return nil, err
		}
		defer os.Remove(file)
	}
	var idFaces []goFace.Face
	if _this.opt.UseCNN {
		idFaces, err = _this.rec.RecognizeFileCNN(file)
	} else {
		idFaces, err = _this.rec.RecognizeFile(file)
	}
	if err != nil {
		return nil, fmt.Errorf("Can't recognize: %v", err)
	}
	return idFaces, nil
}

/*
Classify returns all faces identified in the image. Empty list is returned if no match.
*/
func (_this *Recognizer) Classify(Path string) ([]Face, error) {
	face, err := _this.RecognizeSingle(Path)
	if err != nil {
		return nil, err
	}
	personID := _this.rec.ClassifyThreshold(face.Descriptor, _this.opt.Tolerance)
	if personID < 0 {
		return nil, fmt.Errorf("Can't classify")
	}
	facesRec := make([]Face, 0)
	aux := Face{Data: _this.dataset[personID], Rectangle: face.Rectangle}
	facesRec = append(facesRec, aux)
	return facesRec, nil
}

/*
ClassifyMultiples returns all faces identified in the image. Empty list is returned if no match.
*/
func (_this *Recognizer) ClassifyMultiples(Path string) ([]Face, error) {
	faces, err := _this.RecognizeMultiples(Path)
	if err != nil {
		return nil, fmt.Errorf("Can't recognize: %v", err)
	}
	facesRec := make([]Face, 0)
	for _, f := range faces {
		personID := _this.rec.ClassifyThreshold(f.Descriptor, _this.opt.Tolerance)
		if personID < 0 {
			continue
		}
		aux := Face{Data: _this.dataset[personID], Rectangle: f.Rectangle}
		facesRec = append(facesRec, aux)
	}
	return facesRec, nil
}

func (_this *Recognizer) RecognizeByID(dir, id string) (map[string]Face, error) {
	data := make(map[string]Face)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		faces, err := _this.ClassifyMultiples(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		for _, face := range faces {
			if face.Id == id {
				data[file.Name()] = face
				break
			}
		}
	}
	return data, nil
}

func (_this *Recognizer) FilterImageById(dir, id string) (map[string]image.Image, error) {
	images := make(map[string]image.Image)
	faces, err := _this.RecognizeByID(dir, id)
	if err != nil {
		return nil, err
	}
	for i, face := range faces {
		img, err := _this.DrawFaces(filepath.Join(dir, i), []Face{face})
		if err != nil {
			return nil, err
		}
		images[i] = img
	}
	return images, nil
}

type FacesInImage struct {
	Image string
	File  image.Image
	Faces []string
}

func (_this *Recognizer) FilterImageByFacesInImage(dir, file string) ([]FacesInImage, error) {
	faces, err := _this.ClassifyMultiples(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}
	var data []FacesInImage
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		classifiedFaces, err := _this.ClassifyMultiples(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		var foundFaces []Face
		var faceIds []string
		for _, f := range classifiedFaces {
			for _, face := range faces {
				if face.Id == f.Id {
					foundFaces = append(foundFaces, f)
					faceIds = append(faceIds, f.Id)
				}
			}
		}
		if len(faceIds) > 0 {
			img, err := _this.DrawFaces(filepath.Join(dir, file.Name()), foundFaces)
			if err != nil {
				return nil, err
			}
			data = append(data, FacesInImage{
				Image: file.Name(),
				File:  img,
				Faces: faceIds,
			})
		}
	}
	return data, nil
}

/*
fileExists check se file exist
*/
func fileExists(FileName string) bool {
	file, err := os.Stat(FileName)
	return (err == nil) && !file.IsDir()
}

/*
jsonMarshal Marshal interface to array of byte
*/
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
