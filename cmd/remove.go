package cmd

import (
	"bytes"
	"encoding/json"
	"godot-package-manager/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const REMOVE_CMD_NAME_FLAG = "name"

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a dependency of the project",
	Long: `Removes the dependency of the project. If the project file does not exist it will do nothing.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.SetLogLevel(level)
		util.Trace("Log level set to: " + util.GetLogLevel())
		executeRemoveCommand(cmd, args)
	},
}

func executeRemoveCommand(cmd *cobra.Command, args []string){
	var packagePath string = "." + string(os.PathSeparator) + GODOT_PACKAGE
	var gp, err = util.GetGodotPackage(packagePath)
	if err != nil {
		util.Info("Cannot remove.")
		return
	}
	var name string = GetFlagAsString(cmd, REMOVE_CMD_NAME_FLAG)

	if len(name) == 0 {
		util.Trace("Cannot locate flags, trying to load info from arguments.")
		if len(args) == 0 {
			util.Trace("Cannot load properties from arguments.")
			return
		}
		name = args[0]
	}

	util.Trace("Removing dependency - Name: " + name)
	var plugins []util.GPPlugin = []util.GPPlugin{}
	for _, plugin := range gp.Plugins {
		if plugin.Name != name {
			plugins = append(plugins, plugin)
		}
	}
	if len(plugins) == len(gp.Plugins) {
		util.Info("Dependency not found.")
		return
	}
	gp.Plugins = plugins
	gpBytes := new(bytes.Buffer)
	json.NewEncoder(gpBytes).Encode(gp)
	util.WriteToFile(packagePath, gpBytes.Bytes())
	var splitName []string = strings.Split(name, "/")
	os.RemoveAll("./" + ADDONS + "/" + splitName[len(splitName) - 1])
	util.Info("Dependency removed.")
}

func init() {
	removeCmd.SetUsageTemplate(`Usage:
    gpm remove dependency
    gpm remove --name=Dependency`)

	removeCmd.Flags().String(REMOVE_CMD_NAME_FLAG, "", "Dependency name")

	rootCmd.AddCommand(removeCmd)
}
