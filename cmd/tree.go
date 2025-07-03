/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"
	Logs "dirvcs/internal/services/logging"
	"dirvcs/internal/services/treelogs"
	"fmt"

	"github.com/spf13/cobra"
)

var index int
var list bool
var uuid string
var removeTreeId string

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		Init.CheckInit()

		hasIndex := cmd.Flags().Changed("list")

		if hasIndex {
			treelogs.PrintTreeLogs()
		} else if uuid != "" {
			dirvcs.PrintTreeUUID(uuid)
		} else if removeTreeId != "" {
			treelogs.DeleteLogUuid(removeTreeId)
			Logs.AppendLog(fmt.Sprintf("tree deleted %s", removeTreeId))
		} else {
			dirvcs.PrintTree(index)
		}

	},
}

func init() {
	treeCmd.Flags().IntVarP(&index, "index", "i", 0, "Previous Index of Persisted Tree")
	treeCmd.Flags().BoolVarP(&list, "list", "l", false, "List all persisted trees")
	treeCmd.Flags().StringVar(&uuid, "uuid", "", "List all persisted trees")
	treeCmd.Flags().StringVar(&removeTreeId, "remove", "", "Removes a persisted trees given its UUID")

	rootCmd.AddCommand(treeCmd)
}
