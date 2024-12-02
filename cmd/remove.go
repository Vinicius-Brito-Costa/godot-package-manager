package cmd

import (
	"bytes"
	"encoding/json"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const REMOVE_CMD_NAME_FLAG = "name"

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a dependency of the project",
	Long:  `Removes the dependency of the project. If the project file does not exist it will do nothing.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetLogLevel(level)
		logger.Trace("Log level set to: " + logger.GetLogLevel())
		executeRemoveCommand(cmd, args)
	},
}

func executeRemoveCommand(cmd *cobra.Command, args []string) {
	var packagePath string = "." + string(os.PathSeparator) + GODOT_PACKAGE
	var gp, err = file.GetGodotPackage(packagePath)
	if err != nil {
		logger.Info("Cannot remove.")
		return
	}
	var name string = GetFlagAsString(cmd, REMOVE_CMD_NAME_FLAG)

	if len(name) == 0 {
		logger.Trace("Cannot locate flags, trying to load info from arguments.")
		if len(args) == 0 {
			logger.Trace("Cannot load properties from arguments.")
			return
		}
		name = args[0]
	}

	logger.Trace("Removing dependency - Name: " + name)
	var plugins []file.GPPlugin = []file.GPPlugin{}
	for _, plugin := range gp.Plugins {
		if plugin.Name != name {
			plugins = append(plugins, plugin)
		}
	}
	if len(plugins) == len(gp.Plugins) {
		logger.Info("Dependency not found.")
		return
	}
	gp.Plugins = plugins
	gpBytes := new(bytes.Buffer)
	json.NewEncoder(gpBytes).Encode(gp)
	file.WriteToFile(packagePath, gpBytes.Bytes())
	var splitName []string = strings.Split(name, "/")
	os.RemoveAll("." + string(os.PathSeparator) + ADDONS + string(os.PathSeparator) + splitName[len(splitName)-1])
	logger.Info("Dependency removed.")
}

func init() {
	removeCmd.SetUsageTemplate(`Usage:
    gpm remove dependency
    gpm remove --name=Dependency`)

	removeCmd.Flags().String(REMOVE_CMD_NAME_FLAG, "", "Dependency name")

	rootCmd.AddCommand(removeCmd)
}
