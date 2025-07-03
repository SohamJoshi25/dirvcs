var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Compare two directory snapshots or two exported snapshot files",

	Long: `The 'changes' command compares two directory states using either:
- Snapshot UUIDs stored internally (via --old and --new)
- Or exported .gz snapshot files (via --old-path and --new-path)

If neither --old nor --new is provided, the default comparison is:
  last saved snapshot vs. current working directory.

You can also compare two .gz exported snapshots directly using absolute file paths.`,

	Example: `
  # Compare two snapshots by UUID
  dirvcs changes --old abc-uuid --new def-uuid

  # Compare last snapshot vs current directory
  dirvcs changes --old abc-uuid
  dirvcs changes

  # Compare two exported snapshot files directly
  dirvcs changes --old-path /full/path/1.gz --new-path /full/path/2.gz
`,
	Run: func(cmd *cobra.Command, args []string) {
		// (logic untouched)
	},
}

func init() {
	changesCmd.Flags().StringVarP(&oldId, "old", "o", "", "UUID of the base snapshot (optional, defaults to last snapshot)")
	changesCmd.Flags().StringVarP(&newId, "new", "n", "", "UUID of the target snapshot (optional, defaults to current directory)")

	changesCmd.Flags().StringVar(&oldPath, "old-path", "", "Absolute path to the first .gz snapshot file (used in file comparison mode)")
	changesCmd.Flags().StringVar(&newPath, "new-path", "", "Absolute path to the second .gz snapshot file (used in file comparison mode)")

	changesCmd.Flags().StringVar(&exportJson, "export-verbose", "./changelog.json", "Path to export verbose change logs (JSON format).\nExample: /path/to/changelog.json\nDefaults to './changelog.json'")
	changesCmd.Flags().StringVar(&exportTxt, "export-simple", "./changelog.txt", "Path to export simplified change logs (TXT format).\nExample: /path/to/changelog.txt\nDefaults to './changelog.txt'")

	changesCmd.Flags().BoolVarP(&Print, "print", "p", viper.GetBool("changes.print"), "Print the changelog to the terminal")
	changesCmd.Flags().BoolVarP(&forceverbose, "verbose", "v", viper.GetBool("verbose"), "Force verbose comparison mode")
	changesCmd.Flags().BoolVarP(&forceexport, "export", "e", viper.GetBool("changes.export"), "Force exporting the changelog to file")

	rootCmd.AddCommand(changesCmd)
}
