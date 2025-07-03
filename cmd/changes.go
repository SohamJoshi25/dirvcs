package cmd

import (
	Dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var (
	oldId string
	newId string
)

var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Compare two directory tree states",
	Long: `The 'changes' command compares two tree snapshots by UUID,
or compares a previous snapshot to the current working directory.

If --old is not provided, the most recent snapshot is used as the base.
If --new is not provided, it compares against the current directory.`,
	Example: `
  dirvcs changes --old abc-uuid --new def-uuid
  dirvcs changes --old abc-uuid
  dirvcs changes                     # default: last snapshot vs working directory
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()
		Dirvcs.CompareTree(oldId, newId)
	},
}

func init() {
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of base snapshot (optional, defaults to last)")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of target snapshot (optional, defaults to current working directory)")

	rootCmd.AddCommand(changesCmd)
}
