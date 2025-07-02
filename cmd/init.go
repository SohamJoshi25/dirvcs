/*
Copyright © 2025 Soham Joshi <sohamjoshichinchwad@gmail.com>

*/

package cmd

import (
	createInit "dirvcs/internal/services"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises dirvcs in the current directory.",
	Long: `Initialisess dirvcs in the current working directory. 
		Creates .dirvcs directory to store metadata. 
		Creates .ignore, config.yaml, logs.json, trees/treelogs.json to store metadata about directory.
		Without a .dirvcs folder, dirvcs cannot work.`,
	Run: func(cmd *cobra.Command, args []string) {
		createInit.CreateInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
