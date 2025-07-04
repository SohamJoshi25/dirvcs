/*
Copyright © 2025 Soham Joshi <sohamjoshichinchwad@gmail.com>

*/

package cmd

import (
	Color "dirvcs/internal/services/color"
	Init "dirvcs/internal/services/init"
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises dirvcs in the current directory.",
	Long: `Initialisess dirvcs in the current working directory. 
		Creates .dirvcs directory to store metadata. 
		Creates .ignore, config.yaml, logs.json, trees/treelogs.json to store metadata about directory.
		Without a .dirvcs folder, dirvcs cannot work except when comparing trees with specified .gz path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Color.Color("Welcome to DIR VCS\n", Color.Cyan))
		Init.CreateInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
