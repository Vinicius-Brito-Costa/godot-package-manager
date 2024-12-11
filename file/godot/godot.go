package godot

import (
	"errors"
	"fmt"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"os"
	"strconv"
	"strings"
)

const AUTOLOAD_TAG = "[autoload]"
const EDITOR_PLUGINS_TAG = "[editor_plugins]"
const PLUGIN_TAG = "[plugin]"
const PLUGIN_CFG_FILE = "plugin.cfg"
const GODOT_PROJECT_FILE = "project.godot"
const GODOT_PLUGIN_RES_PATH = "*res://"
const GODOT_PATH_SEPARATOR = "/"
const COMMENT_PREFIX = ";"
const BREAK_LINE = "\n"

const KEY_VALUE string = "key-value"
const TAG string = "tag"
const UNMAPPED string = "unmapped"
const EMPTY string = "empty"
const COMMENT string = "comment"

const PACKED_STRING_ARRAY_START string = "enabled=PackedStringArray("

type PluginConfig struct {
	Name        string
	Description string
	Author      string
	Version     string
	Script      string
	ActivateNow bool
}

type LinkedList struct {
	NextNode     *LinkedList
	PreviousNode *LinkedList
	Data         string
	Metadata     []string
}

func LoadGodotProjectFile() (LinkedList, error) {

	fileData, err := file.GetFile("."+string(os.PathSeparator)+GODOT_PROJECT_FILE, true)
	if err != nil {
		logger.Trace("Cannot load " + GODOT_PROJECT_FILE)
		return LinkedList{}, err
	}

	var mappedValues *LinkedList = new(LinkedList)
	var head = mappedValues
	var currentConfigTag string = ""
	var lineCount int = len(strings.Split(string(fileData), BREAK_LINE))
	logger.Info(strconv.Itoa(lineCount))
	for i, line := range strings.Split(string(fileData), BREAK_LINE) {
		if strings.HasPrefix(line, COMMENT_PREFIX) {
			mappedValues.Metadata = []string{COMMENT}
		} else if isTag(line) && len(strings.TrimSpace(line)) > 2 {
			currentConfigTag = strings.TrimSpace(line)
			mappedValues.Metadata = []string{currentConfigTag, TAG}
		} else if len(currentConfigTag) > 0 && isKeyValue(line) {
			mappedValues.Metadata = []string{currentConfigTag, KEY_VALUE}
		} else if len(currentConfigTag) == 0 {
			mappedValues.Metadata = []string{currentConfigTag, UNMAPPED}
		} else {
			mappedValues.Metadata = []string{currentConfigTag, EMPTY}
		}
		mappedValues.Data = line
		if i < lineCount-1 {
			mappedValues.NextNode = new(LinkedList)
			mappedValues.NextNode.PreviousNode = mappedValues
			mappedValues = mappedValues.NextNode
		}
	}

	return *head, nil
}

func SaveGodotProjectFile(godotProject *LinkedList) bool {
	var data string = ""
	continueLoop := true
	for continueLoop {
		data += godotProject.Data + BREAK_LINE
		if godotProject.NextNode != nil {
			godotProject = godotProject.NextNode
		} else {
			continueLoop = false
		}
	}
	if !file.WriteToFile("." + string(os.PathSeparator) + GODOT_PROJECT_FILE, []byte(data)) {
		logger.Trace("Cannot write updates to file...")
		return false
	}

	logger.Info("Succesfully updated " + GODOT_PROJECT_FILE)

	return true
}

func ActivatePluginOnProject(pluginFolderPath string) bool {
	logger.Info("Setting the plugin up..")
	logger.Trace("Path to plugin: " + pluginFolderPath)
	cfg, err := LoadCFGExtension(pluginFolderPath + string(os.PathSeparator) + PLUGIN_CFG_FILE)
	if err != nil {
		logger.Error("Cannot load "+PLUGIN_CFG_FILE, err)
		return false
	}
	logger.Info("Loaded " + PLUGIN_CFG_FILE)

	head, err := LoadGodotProjectFile()
	if err != nil {
		logger.Error("Cannot load "+GODOT_PROJECT_FILE, err)
		return false
	}
	logger.Info("Loaded " + GODOT_PROJECT_FILE)

	var pluginLine string = "\"" + GODOT_PLUGIN_RES_PATH + strings.ReplaceAll(strings.ReplaceAll(pluginFolderPath, "."+string(os.PathSeparator), ""), string(os.PathSeparator), GODOT_PATH_SEPARATOR) + GODOT_PATH_SEPARATOR + cfg.Script + "\""
	var pluginLineCfg string = replaceScriptWithCfgFromPath(strings.ReplaceAll(pluginLine, "\"*res", "\"res"))
	var continueLoop bool = true
	var root *LinkedList = &head
	var isEditorPluginSet bool = false
	for continueLoop {
		var tag string = root.Metadata[0]
		if !strings.EqualFold(tag, root.Data) {
			if !isEditorPluginSet && strings.EqualFold(EDITOR_PLUGINS_TAG, tag) {
				var valueType string = root.Metadata[1]
				if len(root.Data) > 0 {
					if strings.HasPrefix(root.Data, PACKED_STRING_ARRAY_START) {
						isEditorPluginSet = true
						if strings.Contains(root.Data, pluginLineCfg) {
							logger.Trace("Plugin is already registered on " + EDITOR_PLUGINS_TAG)
							continue
						}
						var currentStringArray []string = getPackedArrayStringContents(root.Data)
						currentStringArray = append(currentStringArray, pluginLineCfg)
						root.Data = createPackedArrayString(currentStringArray)
						logger.Trace("Plugin registered on " + EDITOR_PLUGINS_TAG)
					}
				} else if root.NextNode == nil || valueType == TAG {
					appendPreviousNewNode(root, createPackedArrayString([]string{pluginLineCfg})+BREAK_LINE, []string{EDITOR_PLUGINS_TAG, KEY_VALUE})
					isEditorPluginSet = true
					logger.Trace("Plugin registered on a new line in " + EDITOR_PLUGINS_TAG)
				}
			}
		}

		if root.NextNode != nil {
			root = root.NextNode
		} else {
			continueLoop = false
		}
	}

	if !isEditorPluginSet {
		appendNextNewNode(root, EDITOR_PLUGINS_TAG, []string{EDITOR_PLUGINS_TAG, TAG})
		root = root.NextNode
		appendNextNewNode(root, createPackedArrayString([]string{pluginLineCfg})+BREAK_LINE, []string{EDITOR_PLUGINS_TAG, KEY_VALUE})
		isEditorPluginSet = true
		logger.Trace(EDITOR_PLUGINS_TAG + " tag created and " + pluginLineCfg + " added.")
	}

	return SaveGodotProjectFile(&head)
}
func replaceScriptWithCfgFromPath(path string) string {
	var pluginSplitted []string = strings.Split(path, GODOT_PATH_SEPARATOR)
	return strings.ReplaceAll(path, pluginSplitted[len(pluginSplitted)-1], PLUGIN_CFG_FILE) + "\""
}
func appendNextNewNode(nodes *LinkedList, data string, metadata []string) {
	var newNode *LinkedList = new(LinkedList)
	newNode.Data = data
	newNode.Metadata = metadata
	newNode.PreviousNode = nodes
	nodes.NextNode = newNode
}

func appendPreviousNewNode(nodes *LinkedList, data string, metadata []string) {
	var newNode *LinkedList = new(LinkedList)
	newNode.Data = data
	newNode.Metadata = metadata
	newNode.NextNode = nodes
	newNode.PreviousNode = nodes.PreviousNode
	nodes.PreviousNode = newNode
}

func createPackedArrayString(content []string) string {
	return PACKED_STRING_ARRAY_START + strings.Join(content, ",") + ")"
}
func getPackedArrayStringContents(str string) []string {
	return strings.Split(strings.ReplaceAll(strings.ReplaceAll(str, PACKED_STRING_ARRAY_START, ""), ")", ""), ",")
}
func LoadCFGExtension(path string) (PluginConfig, error) {
	file, err := file.GetFile(path, true)

	if err != nil {
		logger.Error(err.Error(), err)
		return PluginConfig{}, err
	}
	var config PluginConfig
	var hasPluginTag bool = false
	for _, line := range strings.Split(string(file), BREAK_LINE) {

		if isTag(strings.TrimSpace(line)) {
			hasPluginTag = strings.EqualFold(PLUGIN_TAG, strings.TrimSpace(line))
		}

		if len(line) == 0 || strings.HasPrefix(line, COMMENT_PREFIX) {
			continue
		}
		var kv []string = strings.SplitN(line, "=", 2)
		if len(kv) < 2 {
			continue
		}
		kv[0] = strings.TrimSpace(kv[0])
		kv[1] = strings.ReplaceAll(kv[1], "\"", "")
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

	if !hasPluginTag {
		return PluginConfig{}, errors.New("the file does have a plugin tag")
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
