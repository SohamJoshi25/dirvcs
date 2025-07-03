package cmd

import (
	"dirvcs/internal/services/ignore"
	Init "dirvcs/internal/services/init"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listIgnore   bool
	removeIgnore string
)

var ignoreCmd = &cobra.Command{
	Use:   "ignore [patterns]",
	Short: "Manage ignore rules for files and folders",
	Long: `The 'ignore' command lets you manage which files or folders should be ignored by DirVCS.

You can:
- Add patterns to ignore
- List current ignore patterns
- Remove specific patterns from ignore list`,
	Example: `
  dirvcs ignore node_modules build
  dirvcs ignore --list
  dirvcs ignore --remove build
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()

		if removeIgnore != "" {
			ignore.RemoveIgnore(removeIgnore)
			fmt.Println("Pattern removed from ignore list.")
			return
		}

		if listIgnore {
			ignore.PrintIgnore()
			return
		}

		if len(args) == 0 {
			fmt.Println("Error: Please provide at least one pattern or use --list / --remove")
			return
		}

		ignore.ApendIgnore(args)
		fmt.Println("Patterns added to ignore list.")
	},
}

func init() {
	ignoreCmd.Flags().BoolVarP(&listIgnore, "list", "l", false, "List all ignored files and folders")
	ignoreCmd.Flags().StringVar(&removeIgnore, "remove", "", "Remove a pattern from the ignore list")

	rootCmd.AddCommand(ignoreCmd)
}
