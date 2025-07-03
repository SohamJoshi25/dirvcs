package cmd

import (
	"fmt"

	Dirvcs "dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var ExportTreeUUID string
var path string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a directory tree to .json format",
	Long: `Export a directory as a versioned tree.
You can optionally specify a UUID, or it will use default logic.`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()

		if ExportTreeUUID == "" {
			fmt.Println("Exporting directory tree without UUID...")
		} else {
			fmt.Printf("Exporting directory tree with UUID: %s\n", ExportTreeUUID)
		}

		fmt.Printf("Path to export: %s\n", path)

		Dirvcs.ExportTree(ExportTreeUUID, path)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&path, "path", "p", "", "Path to directory (required)")
	exportCmd.Flags().StringVarP(&ExportTreeUUID, "uuid", "u", "", "Optional UUID for export")
	exportCmd.MarkFlagRequired("path")
}
