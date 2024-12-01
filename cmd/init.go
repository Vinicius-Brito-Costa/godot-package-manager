package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"godot-package-manager/util"
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
		util.SetLogLevel(level)
		util.Trace("Log level set to: " + util.GetLogLevel())
		util.Info("Initiating the project...")
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
		util.Trace("Cannot locate flags, trying to load info from arguments.")
		if len(args) < 3 {
			util.Trace("Cannot init project, no arguments provided. Prompting the user...")
			
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
			for index, arg := range args {
				if index == 0 {
					name = arg
				}
				if index == 1 {
					projectType = arg
				}
				if index == 2 {
					version = arg
				}
			}
		}
	}
	var godotPackage util.GodotPackage = util.GodotPackage{}
	godotPackage.Project = util.GPProject{}
	godotPackage.Project.Name = name
	godotPackage.Project.Type = projectType
	godotPackage.Project.Description = ""
	godotPackage.Project.GodotVersion = ""
	godotPackage.Project.Version = version
	godotPackage.Plugins = []util.GPPlugin{}

	var deps []util.GodotPackage = *util.LoadGodotPackagesFromDirectory("./" + ADDONS, GODOT_PACKAGE)
	for _, dep := range deps {
		var gpPlugin util.GPPlugin = util.GPPlugin{}
		gpPlugin.Name = dep.Project.Name
		gpPlugin.Repository = dep.Project.Repository
		gpPlugin.Version = dep.Project.Version
		godotPackage.Plugins = append(godotPackage.Plugins, gpPlugin)
	}

	util.Trace("Creating " + GODOT_PACKAGE + " file with values - Name: " + name + " - Type: " + projectType + " - Version: " + version)

	godotPackageBytes := new(bytes.Buffer)
	json.NewEncoder(godotPackageBytes).Encode(godotPackage)

	godotPackageBytes.Bytes() // this is the []byte
	util.WriteToFile("./"+GODOT_PACKAGE, godotPackageBytes.Bytes())
}

// Get flag as string, if any error occuors it will return an empty string
func GetFlagAsString(cmd *cobra.Command, flagName string) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		util.Trace("Cannot get flag. err: " + err.Error())
		return ""
	}

	return flagValue
}

func stringIsEmpty(str string) bool {
	return len(str) == 0
}

func init() {
	initCmd.SetUsageTemplate(`Usage:
    gpm init`)

	initCmd.Flags().String(INSTALL_CMD_NAME_FLAG, "", "Project name")
	initCmd.Flags().String(INSTALL_CMD_TYPE_FLAG, "", "Project type")
	initCmd.Flags().String(INSTALL_CMD_VERSION_FLAG, "0.0.1", "Project version")

	rootCmd.AddCommand(initCmd)

}
