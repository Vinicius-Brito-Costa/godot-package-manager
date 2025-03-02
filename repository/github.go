package repository

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	copyUtil "godot-package-manager/gpm/copy"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/file/godot"
	"godot-package-manager/gpm/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Github struct{}

type GithubAuthentication struct {
	Token string `json:"token"`
}
type GithubConfiguration struct {
	Authentication GithubAuthentication `json:"authentication"`
}
type UrlTemplate struct {
	Name    string
	Version string
	Package string
}

const API_URL_TEMPLATE string = "http://api.github.com/repos/{{.Name}}/zipball/{{.Version}}"

var URL_TEMPLATES []string = []string{
	"https://github.com/{{.Name}}/archive/refs/tags/{{.Version}}.zip",
	"https://github.com/{{.Name}}/releases/download/{{.Version}}/{{.Package}}-{{.Version}}.zip",
}

func (g Github) Config(plugin *file.GPPlugin) *[]byte {
	// TODO: Get config from global (environment variables, an glboal file with config)
	if plugin == nil || plugin.Config == nil {
		return nil
	}
	var arrObj, err = json.Marshal(plugin.Config)
	if err != nil {
		logger.Error("Cannot marshal plugin config.", err)
		return nil
	}
	if len(arrObj) < 1 {
		logger.Trace("There's no config on plugin")
		return nil
	}

	return &arrObj
}

func (g Github) Download(plugin file.GPPlugin, destiny string) bool {
	if len(plugin.Name) == 0 || len(plugin.Version) == 0 {
		logger.Info("Cannot download. Name or Version missing. Name: " + plugin.Name + " Version: " + plugin.Version)
		return false
	}
	var response *http.Response
	var err error
	var config = g.Config(&plugin)
	var githubConfig GithubConfiguration
	var hasAuth bool = false

	if config != nil {
		err = json.Unmarshal(*config, &githubConfig)
		if err == nil {
			hasAuth = true
		} else {
			logger.Warn("Cannot parse config for " + plugin.Name + " error: " + err.Error())
			// Explicity setting err to nil even when the next code block will override it
			// Why? In the future, the next code block could be changed and
			// the knowledge that the err needs to be null can be forgotten,
			// this can cause various bugs and would be hard to debug it.
			err = nil
		}
	}

	if hasAuth {
		response, err = getWithAuthentication(plugin.Name, plugin.Version, &githubConfig)
		if err != nil {
			logger.Error("Github api call failed.", err)
			err = nil
		}
	}

	// If the authenticated call fails, we will try to get it
	// with an unauthenticated call either way.
	if response == nil {
		response, err = getUntil(plugin.Name, plugin.Version)
	}

	if err != nil {
		logger.Error(err.Error(), err)
		return false
	}

	if response == nil {
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
		if strings.Contains(zipFile.Name, godot.PLUGIN_CFG_FILE) {
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

	if len(folderWithFiles) < 1 {
		logger.Trace("Plugin " + plugin.Name + plugin.Version + " does not have an " + godot.PLUGIN_CFG_FILE + " file.")
		return false
	}
	var split = strings.Split(folderWithFiles, string(os.PathSeparator))
	var current string = destiny + string(os.PathSeparator) + folderWithFiles
	var target string = destiny + string(os.PathSeparator) + split[len(split)-1]
	copyUtil.Dir(current, target)
	os.RemoveAll(destiny + string(os.PathSeparator) + split[0])

	godot.ActivatePluginOnProject(target)
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

func getWithAuthentication(name string, version string, config *GithubConfiguration) (*http.Response, error) {
	var tmpl, err = template.New("api-temp-template-github").Parse(API_URL_TEMPLATE)
	if err != nil {
		logger.Trace("Cannot create api url template. " + API_URL_TEMPLATE)
		return nil, err
	}

	var url bytes.Buffer
	err = tmpl.Execute(&url, UrlTemplate{name, version, name})
	if err != nil {
		logger.Trace("Cannot apply api template (" + url.String() + "). Err: " + err.Error())
		return nil, err
	}

	var client = &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Authorization", "Bearer "+config.Authentication.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("accept", "application/vnd.github+json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Trace("Cannot download with url (" + url.String() + "). Err: " + err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		logger.Warn("GET request on " + url.String() + " failed. Status: " + resp.Status)
		return nil, nil
	}
	err = nil
	logger.Trace("Successfuly downloaded " + name + ":" + version)
	return resp, err
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
		tmplErr = tmpl.Execute(&buff, urlTmpl)
		if tmplErr != nil {
			logger.Trace("Cannot apply template (" + url + "). Err: " + tmplErr.Error())
			responseErr = tmplErr
			continue
		}
		resp, reqErr := http.Get(buff.String())

		if reqErr != nil {
			logger.Trace("Cannot download with url (" + buff.String() + "). Err: " + reqErr.Error())
			responseErr = reqErr
			continue
		}

		if resp.StatusCode != 200 {
			logger.Warn("GET request on " + buff.String() + " failed. Status: " + resp.Status)
			responseErr = errors.New("GET request on " + buff.String() + " failed.")
			continue
		}
		responseErr = nil
		logger.Trace("Successfuly downloaded " + name + ":" + version)
		return resp, responseErr
	}
	return &http.Response{}, responseErr
}
