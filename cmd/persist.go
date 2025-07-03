/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var message string

var persistCmd = &cobra.Command{
	Use:   "persist",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()
		dirvcs.GenerateTree(".", message)
	},
}

func init() {
	rootCmd.AddCommand(persistCmd)
	persistCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message for the snapshot")

}
