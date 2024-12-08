package godot

import (
	"errors"
	"fmt"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"strings"
)

const AUTOLOAD_TAG = "[autoload]"
const PLUGIN_TAG = "[plugin]"
const PLUGIN_CFG_FILE = "plugin.cfg"

type PluginConfig struct {
	Name        string
	Description string
	Author      string
	Version     string
	Script      string
	ActivateNow bool
}

func LoadGodotProjectFile() {

}

func SaveGodotProjectFile(){
	
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
		var kv []string = strings.SplitN(line, "=", 2)
		kv[0] = strings.TrimSpace(kv[0])
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
	logger.Info(strings.Join(splittedKV, ","))
	return false
}
func LoadPluginConfig(path string) {
	return
}
func getAutoloadTagLine() int {

	return -1
}
