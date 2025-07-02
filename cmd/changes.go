/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	Dirvsc "dirvcs/internal/dirvcs"

	"github.com/spf13/cobra"
)

var oldId string
var newId string

// changesCmd represents the changes command
var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Dirvsc.CompareTree(oldId, newId)
	},
}

func init() {
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of persisit which is to be considered base. If not given, will compare the last persist.")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of persisit which is to be compared with. If not given, will compare to current unpersisted working directory")
	rootCmd.AddCommand(changesCmd)
}
