// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package cmd

import (
	Init "dirvcs/internal/services/init"
	Logs "dirvcs/internal/services/logging"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Display internal operation logs for DirVCS",
	Long: `The 'logs' command shows a chronological list of DirVCS internal operations
such as persist actions, deletes, comparisons, and more.

Useful for understanding what actions were performed.`,
	Example: `
  dirvcs logs
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()
		Logs.PrintLogs()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
