package cmd

import (
	"dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var message string

var persistCmd = &cobra.Command{
	Use:   "persist",
	Short: "Snapshot and persist the current directory state",
	Long: `The 'persist' command captures the current state of your directory
and saves it as a versioned snapshot.

You can provide a commit message using the -m flag.`,
	Example: `
  dirvcs persist -m "Initial project snapshot"
  dirvcs persist --message "After refactoring utils"
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()
		dirvcs.GenerateTree(".", message)
	},
}

func init() {
	persistCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message for the snapshot")
	persistCmd.MarkFlagRequired("message")

	rootCmd.AddCommand(persistCmd)
}
