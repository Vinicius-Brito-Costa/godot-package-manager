package cmd

import (
	"godot-package-manager/gpm/file"
	"godot-package-manager/gpm/logger"
	"godot-package-manager/gpm/repository"
	"os"

	"github.com/spf13/cobra"
)

const REPO = "repo: "
const VERSION = "version: "

type Repository interface {
	Download(plugin file.GPPlugin, destiny string) bool
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies",
	Long: `
It uses the ` + GODOT_PACKAGE + ` file inside your project root folder.
The installed plugins will be put in the ` + string(os.PathSeparator) + ADDONS + ` root folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetLogLevel(level)
		logger.Trace("Log level set to: " + logger.GetLogLevel())
		logger.Trace("Installing all dependencies...")
		executeInstallCommand(".")
	},
}

func executeInstallCommand(folder string) {
	var gp, err = file.GetGodotPackage(folder + string(os.PathSeparator) + GODOT_PACKAGE)
	if err != nil {
		logger.Info("Cannot install.")
		return
	}

	logger.Trace("Downloading plugins...")
	for i := range gp.Plugins {
		InstallDependency(gp.Plugins[i])
	}
}
func InstallDependency(dep file.GPPlugin) {
	logger.Info("installing " + dep.Name + ":" + dep.Version)
	var repo Repository = repository.GetRepository(dep.Repository)
	if !repo.Download(dep, "."+string(os.PathSeparator)+ADDONS) {
		var logDir, _ = logger.GetLogDir()
		logger.Info("Cannot download " + dep.Name + ":" + dep.Version + ". Check logs at: " + logDir)
	}
}

// Who knows in the future...
func checkPluginDependencies() {}

func init() {
	installCmd.SetUsageTemplate(`Usage:
    gpm install`)
	rootCmd.AddCommand(installCmd)
}
