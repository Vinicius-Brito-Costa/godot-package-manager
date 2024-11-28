package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)
type GPProject struct {
	Name string;
	Description string;
	Version string;
	GodotVersion string;
}
type GPPlugin struct {
	Repository string;
	Name string;
	Version string;
}
type GodotPackage struct {
	Project GPProject;
	Plugins []GPPlugin;
}

func getFile(path string) (*os.File, error){
	if len(path) == 0 {
		return nil, errors.New("Cannot load file from a blank path.")
	}
	Trace("Getting file in path: " + path)
	file, err := os.Open(path)
	if err != nil {
		Error("Error trying to load file in path: " + path, err)
		return nil, err
	}
	return file, nil
}

func GetGodotPackage(path string) (*GodotPackage) {
	file, err := getFile(path)
	defer file.Close()

	if err != nil {
		Error(err.Error(), err)
		return nil
	}


	r := bufio.NewReader(file)
	var fileData []byte
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
		  fileData = append(fileData, line...)
		}
		if err != nil {
		  break
		}
	}

	if len(fileData) == 0 {
		Info("Cannot load file data.")
		return nil
	}

	var gp GodotPackage

	err = json.Unmarshal(fileData, &gp)

	if err != nil {
		Error(err.Error(), err)
		return nil
	}

	return &gp
}

