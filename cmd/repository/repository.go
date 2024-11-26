package repository

import (
	"archive/zip"
	"bytes"
	"fmt"
	"godot-package-manager/cmd/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const GITHUB string = "github"
const GITLAB string = "gitlab"

type Github struct {}

func (g Github) Download(name string, version string, destiny string) bool {
	if len(name) == 0 || len(version) == 0 {
		util.Info("Cannot download. Name or Version missing. Name: " + name + " Version: " + version)
		return false
	}
	var url = name + "/archive/refs/tags/" + version + ".zip"

	response, err := http.Get(url)

	if err != nil {
		util.Error(err.Error(), err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		util.Info("GET request on " + url + " failed.")
		return false
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

    // Read all the files from zip archive
    for _, zipFile := range zipReader.File {
        err := extract(zipFile, destiny)
        if err != nil {
            util.Error(err.Error(), err)
			continue
        }
    }

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
	if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
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