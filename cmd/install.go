package cmd

import (
	"godot-package-manager/cmd/repository"
	"godot-package-manager/cmd/util"

	"github.com/spf13/cobra"
)

const REPO = "repo: "
const VERSION = "version: "

type Repository interface{
	Download(name string, version string, destiny string) bool
}

// getCmd represents the get command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Info("Installing all dependencies...")
		for _, arg := range args {
			util.Info("Arg: " + arg)
		}
		getGodotPlugins()
	},
}

func getGodotPlugins() {
	var gp util.GodotPackage = *util.GetGodotPackage("./godot-package.json")
	util.Info("Downloading plugins...")
	for i := range gp.Plugins {
		util.Info("Plugin: " + gp.Plugins[i].Repository + " Version: " + gp.Plugins[i].Version)
		var github repository.Github = repository.Github{}
		if github.Download(gp.Plugins[i].Repository, gp.Plugins[i].Version, "./addons") {
			util.Info("Finished downloading " + gp.Plugins[i].Repository)
		}
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
}
