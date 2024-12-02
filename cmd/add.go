package cmd

import (
	"bytes"
	"encoding/json"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"os"

	"github.com/spf13/cobra"
)

const ADD_CMD_NAME_FLAG = "name"
const ADD_CMD_REPOSITORY_FLAG = "repository"
const ADD_CMD_VERSION_FLAG = "version"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a dependency to the project",
	Long:  `Adds the dependency to the project. If the project file does not exist it will do nothing.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetLogLevel(level)
		logger.Trace("Log level set to: " + logger.GetLogLevel())
		executeAddCommand(cmd, args)
	},
}

func executeAddCommand(cmd *cobra.Command, args []string) {
	var packagePath string = "." + string(os.PathSeparator) + GODOT_PACKAGE
	var gp, err = file.GetGodotPackage(packagePath)
	if err != nil {
		logger.Info("Cannot add.")
		return
	}
	var name string = GetFlagAsString(cmd, ADD_CMD_NAME_FLAG)
	var repository string = GetFlagAsString(cmd, ADD_CMD_REPOSITORY_FLAG)
	var version string = GetFlagAsString(cmd, ADD_CMD_VERSION_FLAG)

	if len(name) == 0 || len(repository) == 0 || len(version) == 0 {
		logger.Trace("Cannot locate flags, trying to load info from arguments.")
		if len(args) < 3 {
			logger.Trace("Cannot load properties from arguments.")
			return
		}
		name = args[0]
		repository = args[1]
		version = args[2]
	}

	logger.Trace("Adding dependency - Name: " + name + " - Repository: " + repository + " - Version: " + version)

	var addon file.GPPlugin = file.GPPlugin{}
	addon.Name = name
	addon.Repository = repository
	addon.Version = version
	gp.Plugins = append(gp.Plugins, addon)

	gpBytes := new(bytes.Buffer)
	json.NewEncoder(gpBytes).Encode(gp)
	file.WriteToFile(packagePath, gpBytes.Bytes())
	InstallDependency(addon)
	logger.Info("Dependency added.")
}

func init() {
	addCmd.SetUsageTemplate(`Usage:
    gpm add name repository version
    gpm add --name=Project --repository=github --version=v1.0.0`)

	addCmd.Flags().String(ADD_CMD_NAME_FLAG, "", "Dependency name")
	addCmd.Flags().String(ADD_CMD_REPOSITORY_FLAG, "", "Dependency repository")
	addCmd.Flags().String(ADD_CMD_VERSION_FLAG, "", "Dependency version")

	rootCmd.AddCommand(addCmd)
}
