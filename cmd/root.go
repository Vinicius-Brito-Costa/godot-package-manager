package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gpm",
	Short: "Godot Package Manager",
	Long: `Godot Package Manager (gpm) is a CLI tool for Godot that empowers users to manage plugins.`,
	Version: "0.0.1",

}

func Execute() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true  

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
