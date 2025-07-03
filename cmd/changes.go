package cmd

import (
	Dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var (
	oldId      string
	newId      string
	oldPath    string
	newPath    string
	exportPath string
	skipPrint  bool
)

var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Compare two directory tree states or two exported snapshot files",
	Long: `The 'changes' command compares:
- Two directory tree snapshots by UUID
- A snapshot UUID and current working directory
- Or two .gz exported snapshot files using absolute paths

If --old and --new UUIDs are not provided, it defaults to comparing the last snapshot with the current directory.

Alternatively, you can specify --old-path and --new-path to compare two exported gzip files directly.`,
	Example: `
  dirvcs changes --old abc-uuid --new def-uuid
  dirvcs changes --old abc-uuid
  dirvcs changes
  dirvcs changes --old-path /full/path/1.gz --new-path /full/path/2.gz
`,
	Run: func(cmd *cobra.Command, args []string) {
		if oldPath != "" && newPath != "" {
			Dirvcs.CompareTreePath(oldPath, newPath, exportPath, skipPrint)
			return
		}

		// UUID mode
		Init.CheckInit()
		Dirvcs.CompareTree(oldId, newId, exportPath, skipPrint)
	},
}

func init() {
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of base snapshot (optional, defaults to last)")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of target snapshot (optional, defaults to current working directory)")

	changesCmd.Flags().StringVar(&oldPath, "old-path", "", "Absolute path to first snapshot .gz file")
	changesCmd.Flags().StringVar(&newPath, "new-path", "", "Absolute path to second snapshot .gz file")
	changesCmd.Flags().StringVar(&exportPath, "export", "./changelog.json", "A path to to export change logs to. \ne.g. /path/to/changeslog.tree \nDefaults to './changelog.json' ")
	changesCmd.Flags().BoolVar(&skipPrint, "skip-print", true, "Skips Printing to Command")

	rootCmd.AddCommand(changesCmd)
}
