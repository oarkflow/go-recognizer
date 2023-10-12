package recognizer

import (
	"encoding/json"
	"errors"
	"os"
)

/*
SaveDataset saves dataset data to a json file
*/
func (_this *Recognizer) SaveDataset(Path string) error {
	data, err := jsonMarshal(_this.dataset)
	if err != nil {
		return err
	}
	return os.WriteFile(Path, data, 0777)
}

/*
LoadDataset loads the data from the json file into the dataset
*/
func (_this *Recognizer) LoadDataset(Path string) error {
	if !fileExists(Path) {
		return errors.New("file not found")
	}
	file, err := os.OpenFile(Path, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	Dataset := make([]Data, 0)
	err = json.NewDecoder(file).Decode(&Dataset)
	if err != nil {
		return err
	}
	_this.dataset = append(_this.dataset, Dataset...)
	return nil
}
