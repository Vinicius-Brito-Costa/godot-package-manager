package cmd

import (
	"godot-package-manager/cmd/repository"
	"godot-package-manager/cmd/util"
	"github.com/spf13/cobra"
)

const REPO = "repo: "
const VERSION = "version: "

type Repository interface {
	Download(name string, version string, destiny string) bool
}

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
		getGodotPlugins(".")
	},
}

func getGodotPlugins(folder string) {
	var gp util.GodotPackage = *util.GetGodotPackage(folder + "/godot-package.json")
	util.Trace("Downloading plugins...")
	for i := range gp.Plugins {
		util.Info(gp.Plugins[i].Name + ":" + gp.Plugins[i].Version)
		var repo Repository = repository.GetRepository(gp.Plugins[i].Repository)
		if !repo.Download(gp.Plugins[i].Name, gp.Plugins[i].Version, "./addons") {
			util.Info("Cannot download " + gp.Plugins[i].Name + ":" + gp.Plugins[i].Version)
		}
	}
}

func checkPluginDependencies(){}

func init() {
	rootCmd.AddCommand(installCmd)
}
