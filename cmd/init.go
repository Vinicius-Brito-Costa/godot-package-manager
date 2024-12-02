package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Will start the project",
	Long: `Prompt for some information about the project, it will ask for the name, type and version.
Next it search for dependencies on the addons folder.
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetLogLevel(level)
		logger.Trace("Log level set to: " + logger.GetLogLevel())
		logger.Trace("Initiating the project...")
		executeInitCommand(cmd, args)
	},
}

const INSTALL_CMD_NAME_FLAG = "name"
const INSTALL_CMD_TYPE_FLAG = "type"
const INSTALL_CMD_VERSION_FLAG = "version"

func executeInitCommand(cmd *cobra.Command, args []string) {
	var name string = GetFlagAsString(cmd, INSTALL_CMD_NAME_FLAG)
	var projectType string = GetFlagAsString(cmd, INSTALL_CMD_TYPE_FLAG)
	var version string = GetFlagAsString(cmd, INSTALL_CMD_VERSION_FLAG)

	if len(name) == 0 || len(projectType) == 0 || len(version) == 0 {
		logger.Trace("Cannot locate flags, trying to load info from arguments.")
		if len(args) < 3 {
			logger.Trace("Cannot init project, no arguments provided. Prompting the user...")

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Project name: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Print("Project type: ")
			reader = bufio.NewReader(os.Stdin)
			projectType, _ = reader.ReadString('\n')
			projectType = strings.TrimSpace(projectType)
			fmt.Print("Project version: ")
			reader = bufio.NewReader(os.Stdin)
			version, _ = reader.ReadString('\n')
			version = strings.TrimSpace(version)
		} else {
			name = args[0]
			projectType = args[1]
			version = args[2]
		}
	}
	var godotPackage file.GodotPackage = file.GodotPackage{}
	godotPackage.Project = file.GPProject{}
	godotPackage.Project.Name = name
	godotPackage.Project.Type = projectType
	godotPackage.Project.Description = ""
	godotPackage.Project.GodotVersion = ""
	godotPackage.Project.Version = version
	godotPackage.Plugins = []file.GPPlugin{}

	var deps []file.GodotPackage = *file.LoadGodotPackagesFromDirectory("." + string(os.PathSeparator) + ADDONS, GODOT_PACKAGE)
	for _, dep := range deps {
		var gpPlugin file.GPPlugin = file.GPPlugin{}
		gpPlugin.Name = dep.Project.Name
		gpPlugin.Repository = dep.Project.Repository
		gpPlugin.Version = dep.Project.Version
		godotPackage.Plugins = append(godotPackage.Plugins, gpPlugin)
	}

	logger.Trace("Creating " + GODOT_PACKAGE + " file with values - Name: " + name + " - Type: " + projectType + " - Version: " + version)

	godotPackageBytes := new(bytes.Buffer)
	json.NewEncoder(godotPackageBytes).Encode(godotPackage)
	file.WriteToFile("." + string(os.PathSeparator) + GODOT_PACKAGE, godotPackageBytes.Bytes())
	logger.Info(GODOT_PACKAGE + " created.")
}

// Get flag as string, if any error occuors it will return an empty string
func GetFlagAsString(cmd *cobra.Command, flagName string) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		logger.Trace("Cannot get flag. err: " + err.Error())
		return ""
	}

	return flagValue
}

func stringIsEmpty(str string) bool {
	return len(str) == 0
}

func init() {
	initCmd.SetUsageTemplate(`Usage:
    gpm init project type version
    gpm init --name=project --type=game --version=1.0.0`)

	initCmd.Flags().String(INSTALL_CMD_NAME_FLAG, "", "Project name")
	initCmd.Flags().String(INSTALL_CMD_TYPE_FLAG, "", "Project type")
	initCmd.Flags().String(INSTALL_CMD_VERSION_FLAG, "0.0.1", "Project version")

	rootCmd.AddCommand(initCmd)

}
