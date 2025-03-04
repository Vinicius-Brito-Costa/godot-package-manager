package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

const GODOT_PACKAGE = "godot-package.json"
const ADDONS = "addons"
const LOG_LEVEL_FLAG = "log-level"
const logLevelDescription = `Changes log level. Available:
	info	default
	warn
	trace
`

var rootCmd = &cobra.Command{
	Use:     "gpm",
	Short:   "Godot Package Manager",
	Long:    `Godot Package Manager (gpm) is a CLI tool for Godot that empowers users to manage plugins.`,
	Version: "1.1.1",
}
var level string

func Execute() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&level, LOG_LEVEL_FLAG, "", logLevelDescription)

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
