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
			log.Fatalln("Cannot Give Both Parameter 'export-verbose' and 'export-simple'. Please give one")
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
				fmt.Println("Exporting and Printing both cannot be set to false.")
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
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of base snapshot (optional, defaults to last)")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of target snapshot (optional, defaults to current working directory)")

	changesCmd.Flags().StringVar(&oldPath, "old-path", "", "Absolute path to first snapshot .gz file")
	changesCmd.Flags().StringVar(&newPath, "new-path", "", "Absolute path to second snapshot .gz file")
	changesCmd.Flags().StringVar(&exportJson, "export-verbose", "./changelog.json", "A path to to export verbose change logs to. \ne.g. /path/to/changeslog.json \nDefaults to './changelog.json' ")
	changesCmd.Flags().StringVar(&exportTxt, "export-simple", "./changelog.txt", "A path to to export cimple change logs to. \ne.g. /path/to/changeslog.txt \nDefaults to './changelog.txt' ")
	changesCmd.Flags().BoolVarP(&Print, "print", "p", viper.GetBool("changes.print"), "To print the changlog")
	changesCmd.Flags().BoolVarP(&forceverbose, "verbose", "v", viper.GetBool("verbose"), "To override verbose")
	changesCmd.Flags().BoolVarP(&forceexport, "export", "e", viper.GetBool("changes.export"), "To override export")

	rootCmd.AddCommand(changesCmd)
}
