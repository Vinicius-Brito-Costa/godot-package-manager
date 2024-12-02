package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"godot-package-manager/gpm/logger"
	"io/fs"
	"os"
)

type GPProject struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Version      string `json:"version"`
	Repository   string `json:"repository"`
	Description  string `json:"description"`
	GodotVersion string `json:"godotVersion"`
}
type GPPlugin struct {
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}
type GodotPackage struct {
	Project GPProject  `json:"project"`
	Plugins []GPPlugin `json:"plugins"`
}

func getFile(path string) (*os.File, error) {
	if len(path) == 0 {
		return nil, errors.New("Cannot load file from a blank path.")
	}
	logger.Trace("Getting file in path: " + path)
	file, err := os.Open(path)
	if err != nil {
		logger.Error("Error trying to load file in path: "+path, err)
		return nil, err
	}
	return file, nil
}

func GetGodotPackage(path string) (*GodotPackage, error) {
	file, err := getFile(path)
	defer file.Close()

	if err != nil {
		logger.Error(err.Error(), err)
		return &GodotPackage{}, err
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
		logger.Info("Cannot load file data.")
		return &GodotPackage{}, errors.New("Cannot load file data.")
	}

	var gp GodotPackage

	err = json.Unmarshal(fileData, &gp)

	if err != nil {
		logger.Error(err.Error(), err)
		return &GodotPackage{}, err
	}

	return &gp, nil
}

func LoadGodotPackagesFromDirectory(dir string, godotPackageName string) *[]GodotPackage {
	files, err := fs.Glob(os.DirFS(dir), "**" + string(os.PathSeparator) + godotPackageName)
	if err != nil {
		return &[]GodotPackage{}
	}

	var pluginsGodotPackage []GodotPackage = []GodotPackage{}
	for _, file := range files {
		logger.Trace("File: " + file)
		gp, err := GetGodotPackage(dir + string(os.PathSeparator) + file)
		if err != nil {
			continue
		}
		pluginsGodotPackage = append(pluginsGodotPackage, *gp)
	}

	logger.Trace("Number of plugins: " + fmt.Sprint(len(pluginsGodotPackage)))

	return &pluginsGodotPackage
}

func WriteToFile(path string, data []byte) bool {

	var err = os.WriteFile(path, data, 0644)
	if err != nil {
		logger.Error(err.Error(), err)
		return false
	}

	return true
}
