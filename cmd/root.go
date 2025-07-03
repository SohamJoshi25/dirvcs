package cmd

import (
	"fmt"
	"os"

	Path "dirvcs/internal/data/path"
	Color "dirvcs/internal/services/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dirvcs",
	Short: "DirVCS is a lightweight directory version control system",
	Long: `DirVCS is a CLI tool that allows you to snapshot, compare, and manage versions of directory structures.

It supports persisting directory states, comparing changes, pruning old versions, and more.
`,
	Version: "v0.1.0",
	Example: `
  dirvcs init                          # Initialize a new dirvcs repo
  dirvcs persist -m "Snapshot msg"    # Persist current state
  dirvcs changes               # Compare tree with previous snapshot
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Color.Color("Welcome to DIR VCS\n", Color.Blue))
		_ = cmd.Help() // show help by default
	},
}

func init() {
	// Version formatting
	rootCmd.SetVersionTemplate(`{{.Version}}` + "\n")

	// Config file setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(Path.BASE_PATH)
	viper.SetDefault("treelimit", 20)
	viper.SetDefault("verbose", false)
	viper.SetDefault("changes.export", true)
	viper.SetDefault("indent", "|---")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
	}

	// Optional verbose logging
	if viper.GetBool("verbose") {
		fmt.Println("Using config:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
