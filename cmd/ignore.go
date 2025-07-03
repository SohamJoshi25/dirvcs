/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dirvcs/internal/services/ignore"
	Init "dirvcs/internal/services/init"
	"fmt"

	"github.com/spf13/cobra"
)

var listIgnore bool
var removeIgnore string

var ignoreCmd = &cobra.Command{
	Use:   "ignore [patterns]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()

		showList := cmd.Flags().Changed("list")

		if removeIgnore != "" {
			ignore.RemoveIgnore(removeIgnore)
			fmt.Println("Ignore List updated")
		} else if showList {
			ignore.PrintIgnore()
		} else {
			if len(args) == 0 {
				fmt.Println("Please provide at least one pattern to ignore")
				return
			}
			ignore.ApendIgnore(args)
			fmt.Println("Ignore List updated")
		}
	},
}

func init() {
	ignoreCmd.Flags().BoolVarP(&listIgnore, "list", "l", false, "Used to show all files and folders ignored")
	ignoreCmd.Flags().StringVar(&removeIgnore, "remove", "", "Remove from .ignore file")
	rootCmd.AddCommand(ignoreCmd)
}
