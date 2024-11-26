package cmd

import (
	"fmt"
	"godot-package-manager/cmd/util"
	"strings"

	"github.com/spf13/cobra"
)

const REPO = "repo: "
const VERSION = "version: "

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called.")
		for _, arg := range args {
			fmt.Println("Arg: " + arg)
		}
		getGodotPackages()
	},
}

func getGodotPackages() {
	var packageFile, err = util.GetFile("godot-package.txt")
	if err != nil {
		util.Error("File does not exists.", err)
	}

	var lines, err2 = util.GetFileLines(packageFile)
	if err2 != nil {
		util.Error("Cannot get file lines.", err2)
	}
	var dependencies = make(map[string]string)
	for _, line := range lines {
		if !strings.Contains(line, REPO) || !strings.Contains(line, VERSION) {
			util.Info("Cannot read line, invalid format. Line: " + line)
		} else {
			var cleanLine = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "repo:", ""), "version:", ""))
			var infoMap = strings.Fields(cleanLine)
			if len(dependencies[infoMap[0]]) > 0 {
				util.Info("Checking if new dependency has a higher version than the old one.")
			}

			dependencies[infoMap[0]] = infoMap[1]
		}
	}

	for k, v := range dependencies {
		util.Info("Downloading dependencies. Dependency: " + k + " Version: " + v)
	}
	//var url = "/releases/tag/"
}

func init() {
	rootCmd.AddCommand(getCmd)
}
