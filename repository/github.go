package repository

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"godot-package-manager/gpm/logger"
	copyUtil "godot-package-manager/gpm/copy"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Github struct{}

func (g Github) Download(name string, version string, destiny string) bool {
	if len(name) == 0 || len(version) == 0 {
		logger.Info("Cannot download. Name or Version missing. Name: " + name + " Version: " + version)
		return false
	}

	var response, err = getUntil(name, version)

	if err != nil {
		logger.Error(err.Error(), err)
		return false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err.Error(), err)
		return false
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		logger.Error(err.Error(), err)
		return false
	}

	var folderWithFiles string = ""

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
			logger.Error(err.Error(), err)
			continue
		}
	}

	logger.Trace("Folder with files: " + folderWithFiles)

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

type UrlTemplate struct {
	Name    string
	Version string
	Package string
}

var URL_TEMPLATES []string = []string{
	"https://github.com/{{.Name}}/archive/refs/tags/{{.Version}}.zip",
	"https://github.com/{{.Name}}/releases/download/{{.Version}}/{{.Package}}-{{.Version}}.zip",
}

// This will loop over the URL_TEMPLATES looking for a positive. If it cannot find anything, will return error.
func getUntil(name string, version string) (*http.Response, error) {
	// Dinamically getting package name using the name. Don't know if this is the best choice.
	var splitedName = strings.Split(name, "/")
	var urlTmpl UrlTemplate = UrlTemplate{name, version, splitedName[len(splitedName)-1]}
	var responseErr error
	for index, url := range URL_TEMPLATES {
		var tmpl, tmplErr = template.New("temp-tmpl-" + string(index)).Parse(url)
		if tmplErr != nil {
			logger.Trace("Cannot parse template (" + url + "). Err: " + tmplErr.Error())
			responseErr = tmplErr
			continue
		}

		var buff bytes.Buffer
		tmpl.Execute(&buff, urlTmpl)
		resp, reqErr := http.Get(buff.String())

		if reqErr != nil {
			logger.Trace("Cannot download with url (" + buff.String() + "). Err: " + reqErr.Error())
			responseErr = reqErr
			continue
		}

		if resp.StatusCode != 200 {
			logger.Warn("GET request on " + buff.String() + " failed. Status: " + resp.Status)
			if index < len(URL_TEMPLATES) {
				logger.Trace("Trying to GET on the next template. " + URL_TEMPLATES[index + 1])
			}
			responseErr = errors.New("GET request on " + buff.String() + " failed.")
			continue
		}
		responseErr = nil
		logger.Trace("Successfuly downloaded " + name + ":" + version)
		return resp, responseErr
	}
	return &http.Response{}, responseErr
}
