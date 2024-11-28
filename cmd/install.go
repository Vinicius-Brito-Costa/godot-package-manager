package cmd

import (
	"godot-package-manager/repository"
	"godot-package-manager/util"
	"os"

	"github.com/spf13/cobra"
)

const REPO = "repo: "
const VERSION = "version: "
const GODOT_PACKAGE = "godot-package.json"
const ADDONS = "addons"
type Repository interface {
	Download(name string, version string, destiny string) bool
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies",
	Long: `
It uses the ` + GODOT_PACKAGE + ` file inside your project root folder.
The installed plugins will be put in the ` + string(os.PathSeparator) + ADDONS + ` root folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Info("Installing all dependencies...")
		getGodotPlugins(".")
	},
}

func getGodotPlugins(folder string) {
	var gp util.GodotPackage = *util.GetGodotPackage(folder + string(os.PathSeparator) + GODOT_PACKAGE)
	util.Trace("Downloading plugins...")
	for i := range gp.Plugins {
		util.Info(gp.Plugins[i].Name + ":" + gp.Plugins[i].Version)
		var repo Repository = repository.GetRepository(gp.Plugins[i].Repository)
		if !repo.Download(gp.Plugins[i].Name, gp.Plugins[i].Version, "." + string(os.PathSeparator) + ADDONS) {
			util.Info("Cannot download " + gp.Plugins[i].Name + ":" + gp.Plugins[i].Version)
		}
	}
}

func checkPluginDependencies(){}

func init() {
	installCmd.SetUsageTemplate(`Usage:
    gpm install`)
	rootCmd.AddCommand(installCmd)
}
