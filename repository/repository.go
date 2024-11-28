package repository

import (
	"archive/zip"
	"bytes"
	"fmt"
	copyUtil "godot-package-manager/copy"
	"godot-package-manager/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var repositories = map[string]Repository{
	"github": Github{},
}

type Repository interface {
	Download(name string, version string, destiny string) bool
}

type Github struct{}

func (g Github) Download(name string, version string, destiny string) bool {
	if len(name) == 0 || len(version) == 0 {
		util.Info("Cannot download. Name or Version missing. Name: " + name + " Version: " + version)
		return false
	}
	var url = "https://github.com/" + name + "/archive/refs/tags/" + version + ".zip"

	var response, err = http.Get(url)

	if err != nil {
		util.Error(err.Error(), err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		//https://github.com/Burloe/GoLogger/releases/download/1.2/GoLogger-1.2.zip
		var splitedName = strings.Split(name, "/")
		url = "https://github.com/" + name + "/releases/download/" + version + "/" + splitedName[len(splitedName)-1] + "-" + version + ".zip"
		response, err = http.Get(url)
		if err != nil {
			util.Error(err.Error(), err)
			return false
		}
		if response.StatusCode != 200 {
			util.Info("GET request on " + url + " failed.")
			return false
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		util.Error(err.Error(), err)
		return false
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		util.Error(err.Error(), err)
		return false
	}

	var folderWithFiles string = ""
	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {

		// Search for the folder that contains the plugin.cfg file
		if strings.Contains(zipFile.Name, "plugin.cfg") {
			if len(folderWithFiles) == 0 {
				folderWithFiles = filepath.Dir(zipFile.Name)
			}
			// If there's another plugin.cfg file in the project
			// we need to get the upper most folder of the project
			// so we compare the size of the path between them
			// the lowest one wins.
			if len(folderWithFiles) > len(filepath.Dir(zipFile.Name)) {
				folderWithFiles = filepath.Dir(zipFile.Name)
			}
		}

		err := extract(zipFile, destiny)
		if err != nil {
			util.Error(err.Error(), err)
			continue
		}
	}

	util.Trace("Folder with files: " + folderWithFiles)

	var split = strings.Split(folderWithFiles, string(os.PathSeparator))
	copyUtil.Dir(destiny+string(os.PathSeparator)+folderWithFiles, destiny+string(os.PathSeparator)+split[len(split)-1])
	os.RemoveAll(destiny + string(os.PathSeparator) + split[0])
	return true
}

func extract(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(dest, f.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("Illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), f.Mode())
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRepository(repo string) Repository {
	return repositories[repo]
}
