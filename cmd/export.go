// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package cmd

import (
	"dirvcs/internal/dirvcs"
	Init "dirvcs/internal/services/init"

	"github.com/spf13/cobra"
)

var (
	exportTreeUUID string
	path           string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a persisted directory tree to JSON format",
	Long: `The 'export' command outputs a directory tree (by UUID or default) 
into a JSON file.

You must provide the directory path to export the snapshot to.`,
	Example: `
  dirvcs export --path ./output
  dirvcs export --uuid 123e4567-e89b-12d3-a456-426614174000 --path ./tree.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()
		dirvcs.ExportTree(exportTreeUUID, path)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&path, "path", "p", "", "Path to export the tree as JSON (required)")
	exportCmd.Flags().StringVarP(&exportTreeUUID, "uuid", "u", "", "UUID of the tree to export (optional)")
	exportCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(exportCmd)
}
