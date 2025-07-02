package cmd

import (
	Path "dirvcs/internal/data/path"
	color "dirvcs/internal/services/color"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dirvcs",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
	Version: "v0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.Color("Welcome to DIR VCS", color.Blue))
	},
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Version}}` + "\n")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(Path.BASE_PATH)

	viper.SetDefault("prune_limit", 50)
	viper.SetDefault("auto_compress", true)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
