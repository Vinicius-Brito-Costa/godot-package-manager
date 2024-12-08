package godot

import (
	"errors"
	"fmt"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"os"
	"strings"
)

const AUTOLOAD_TAG = "[autoload]"
const PLUGIN_TAG = "[plugin]"
const PLUGIN_CFG_FILE = "plugin.cfg"
const GODOT_PROJECT_FILE = "project.godot"

type PluginConfig struct {
	Name        string
	Description string
	Author      string
	Version     string
	Script      string
	ActivateNow bool
}

func LoadGodotProjectFile() (map[string]map[string]string, error) {

	fileData, err := file.GetFile("."+string(os.PathSeparator)+GODOT_PROJECT_FILE, true)
	if err != nil {
		logger.Trace("Cannot load " + GODOT_PROJECT_FILE)
		return nil, err
	}

	var mappedValues map[string]map[string]string = make(map[string]map[string]string)
	var currentConfigTag string = ""
	for _, line := range strings.Split(string(fileData), "\n") {
		if strings.HasPrefix(line, ";"){
			continue
		}
		if isTag(line) && len(strings.TrimSpace(line)) > 2 {
			currentConfigTag = line
		}
		if len(currentConfigTag) > 0 && isKeyValue(line) {
			var kv []string = strings.SplitN(line, "=", 2)
			if len(kv) > 1 {
				if mappedValues[currentConfigTag] == nil {
					mappedValues[currentConfigTag] = make(map[string]string)
				}
				mappedValues[currentConfigTag][kv[0]] = kv[1]
			}
		}
	}

	return mappedValues, nil
}

func SaveGodotProjectFile(godotProject map[string]map[string]string) {

}
func LoadCFGExtension(path string) (PluginConfig, error) {
	file, err := file.GetFile(path, true)

	if err != nil {
		logger.Error(err.Error(), err)
		return PluginConfig{}, err
	}
	var config PluginConfig
	var hasPluginTag bool = false
	for index, line := range strings.Split(string(file), "\n") {
		if index == 0 {
			hasPluginTag = PLUGIN_TAG == strings.TrimSpace(line)
		}
		if !hasPluginTag {
			return PluginConfig{}, errors.New("The file does not start with a plugin tag")
		}
		if strings.HasPrefix(line, ";"){
			continue
		}
		var kv []string = strings.SplitN(line, "=", 2)
		if len(kv) < 2 {
			continue
		}
		kv[0] = strings.TrimSpace(kv[0])
		logger.Trace("Key: " + kv[0] + " -  Value: " + kv[1])
		if strings.EqualFold(kv[0], "Name") && len(kv[1]) > 0 {
			config.Name = kv[1]
		}
		if strings.EqualFold(kv[0], "Description") && len(kv[1]) > 0 {
			config.Description = kv[1]
		}
		if strings.EqualFold(kv[0], "Author") && len(kv[1]) > 0 {
			config.Author = kv[1]
		}
		if strings.EqualFold(kv[0], "Version") && len(kv[1]) > 0 {
			config.Version = kv[1]
		}
		if strings.EqualFold(kv[0], "Script") && len(kv[1]) > 0 {
			config.Script = kv[1]
		}
		if strings.EqualFold(kv[0], "activate_now") && len(kv[1]) > 0 {
			config.ActivateNow = false
			if strings.Contains(kv[1], "true") {
				config.ActivateNow = true
			}
		}
	}

	logger.Trace("Parsed cfg file: " + fmt.Sprintf("%#v", config))
	return config, nil
}

func isTag(str string) bool {
	if len(str) == 0 {
		return false
	}
	var hasOpeningAndClosingBrackets bool = strings.Index(str, "[") == 0 && strings.Index(str, "]") == len(str)-1

	return hasOpeningAndClosingBrackets
}
func isKeyValue(str string) bool {
	if len(str) == 0 {
		return false
	}
	var splittedKV []string = strings.SplitN(str, "=", 2)
	return len(splittedKV) == 2
}
func LoadPluginConfig(path string) {
	return
}
func getAutoloadTagLine() int {

	return -1
}
