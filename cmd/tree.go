// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package cmd

import (
	"dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"
	Logs "dirvcs/internal/services/logging"
	"dirvcs/internal/services/treelogs"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	index        int
	list         bool
	uuid         string
	removeTreeId string
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "View or manage persisted tree snapshots",
	Long: `The 'tree' command allows you to list, view, or remove persisted tree snapshots.

You can:
- View a tree snapshot by index (--index)
- List all available trees (--list)
- Show a specific tree using its UUID (--uuid)
- Remove a tree by UUID (--remove)
`,
	Example: `
  dirvcs tree --list
  dirvcs tree --index 2
  dirvcs tree --uuid 31abcf99-1234
  dirvcs tree --remove 31abcf99-1234
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()

		if list {
			treelogs.PrintTreeLogs()
			return
		}

		if uuid != "" && removeTreeId != "" {
			fmt.Println("Error: --uuid and --remove cannot be used together")
			return
		}

		if uuid != "" {
			dirvcs.PrintTreeUUID(uuid)
			return
		}

		if removeTreeId != "" {
			treelogs.DeleteLogUuid(removeTreeId)
			fmt.Println("Tree Deleted")
			Logs.AppendLog(fmt.Sprintf("tree deleted %s", removeTreeId))
			return
		}

		// Default to printing by index
		dirvcs.PrintTree(index)
	},
}

func init() {
	treeCmd.Flags().IntVarP(&index, "index", "i", 0, "Index of persisted tree to view")
	treeCmd.Flags().BoolVarP(&list, "list", "l", false, "List all persisted trees")
	treeCmd.Flags().StringVar(&uuid, "uuid", "", "View tree snapshot by UUID")
	treeCmd.Flags().StringVar(&removeTreeId, "remove", "", "Remove persisted tree snapshot by UUID")

	rootCmd.AddCommand(treeCmd)
}
