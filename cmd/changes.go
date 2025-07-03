package cmd

import (
	Dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	oldId        string
	newId        string
	oldPath      string
	newPath      string
	exportJson   string
	exportTxt    string
	Print        bool
	forceverbose bool
	forceexport  bool
)

var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Compare two directory snapshots or exported snapshot files",

	Long: `The 'changes' command compares two directory states using either:

• Snapshot UUIDs stored by dirvcs (via --old and --new)
• Exported snapshot files (.gz) using absolute paths (via --old-path and --new-path)

If neither --old nor --new is provided, the comparison defaults to:
  the last saved snapshot vs. the current working directory.

To compare exported snapshots directly, provide both --old-path and --new-path.
You can choose to print the output, export it to file, or both.`,

	Example: `
  # Compare two snapshots by UUID
  dirvcs changes --old abc-uuid --new def-uuid

  # Compare the last snapshot with the current working directory
  dirvcs changes
  dirvcs changes --old abc-uuid

  # Compare two exported .gz snapshot files
  dirvcs changes --old-path /full/path/1.gz --new-path /full/path/2.gz

  # Export verbose changelog to JSON
  dirvcs changes --export-verbose

  # Export simple changelog to TXT
  dirvcs changes --export-simple

  # Print only, no export
  dirvcs changes --print
`,
	Run: func(cmd *cobra.Command, args []string) {

		currentVerbose := viper.GetBool("verbose")
		currentExport := viper.GetBool("changes.export")

		resetConfig := func() {
			viper.Set("verbose", currentVerbose)
			viper.Set("changes.export", currentExport)
		}

		if cmd.Flags().Changed("verbose") {
			if forceverbose != currentVerbose {
				viper.Set("verbose", forceverbose)
			}
		}
		if cmd.Flags().Changed("export") {
			if forceexport != currentExport {
				viper.Set("changes.export", forceexport)
			}
		}
		defer resetConfig()

		if cmd.Flags().Changed("export-verbose") && cmd.Flags().Changed("export-simple") {
			log.Fatalln("Cannot give both parameters: 'export-verbose' and 'export-simple'. Please provide only one.")
		}

		isUUIDMode := !(cmd.Flags().Changed("old-path") && cmd.Flags().Changed("new-path"))

		verbose := func() {
			if isUUIDMode {
				Init.CheckInit()
				Dirvcs.CompareTree(oldId, newId, exportJson, 2, Print)
			} else {
				Dirvcs.CompareTreePath(oldPath, newPath, exportTxt, 2, Print)
			}
		}

		simple := func() {
			if isUUIDMode {
				Init.CheckInit()
				Dirvcs.CompareTree(oldId, newId, exportTxt, 1, Print)
			} else {
				Dirvcs.CompareTreePath(oldPath, newPath, exportTxt, 1, Print)
			}
		}

		if cmd.Flags().Changed("export-verbose") {
			verbose()
			return
		}

		if cmd.Flags().Changed("export-simple") {
			simple()
			return
		}

		if viper.GetBool("changes.export") {
			if viper.GetBool("verbose") {
				verbose()
			} else {
				simple()
			}
		} else {
			if !Print {
				fmt.Println("Exporting and printing both cannot be set to false.")
			}
			if isUUIDMode {
				Init.CheckInit()
				Dirvcs.CompareTree(oldId, newId, "", 0, true)
			} else {
				Dirvcs.CompareTreePath(oldPath, newPath, "", 0, true)
			}
		}

	},
}

func init() {
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of the base snapshot (optional, defaults to last snapshot)")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of the target snapshot (optional, defaults to current directory)")

	changesCmd.Flags().StringVar(&oldPath, "old-path", "", "Absolute path to the first .gz snapshot file (used in file comparison mode)")
	changesCmd.Flags().StringVar(&newPath, "new-path", "", "Absolute path to the second .gz snapshot file (used in file comparison mode)")

	changesCmd.Flags().StringVar(&exportJson, "export-verbose", "./changelog.json", "Path to export verbose changelog (JSON format).\nExample: --export-verbose /path/to/log.json")
	changesCmd.Flags().StringVar(&exportTxt, "export-simple", "./changelog.txt", "Path to export simple changelog (TXT format).\nExample: --export-simple /path/to/log.txt")

	changesCmd.Flags().BoolVarP(&Print, "print", "p", viper.GetBool("changes.print"), "Print the changelog to the terminal")
	changesCmd.Flags().BoolVarP(&forceverbose, "verbose", "v", viper.GetBool("verbose"), "Force verbose mode (overrides config)")
	changesCmd.Flags().BoolVarP(&forceexport, "export", "e", viper.GetBool("changes.export"), "Force exporting the changelog to file (overrides config)")

	rootCmd.AddCommand(changesCmd)
}
