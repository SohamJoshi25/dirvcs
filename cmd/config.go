package cmd

import (
	"fmt"
	"strconv"

	Init "dirvcs/internal/services/init"
	Log "dirvcs/internal/services/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	setKey   string
	setValue string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update DirVCS configuration",
	Long: `The 'config' command allows you to view current configuration values 
or update specific keys like 'treelimit'.

Settings are stored in the DirVCS YAML config file.`,
	Example: `
  dirvcs config                         # Show all current configuration
  dirvcs config --set-key treelimit --set-value 50
  dirvcs config --set-key autoCompress --set-value true
`,
	Run: func(cmd *cobra.Command, args []string) {
		Init.CheckInit()

		// Update config
		if setKey != "" && setValue != "" {
			if setKey == "treelimit" {
				val, err := strconv.Atoi(setValue)
				if err != nil {
					fmt.Println("Error: treelimit must be an integer")
					return
				}
				if val < 2 {
					fmt.Println("Error: treelimit must be greater than 1")
					return
				}
				viper.Set(setKey, val)
			} else if setKey == "verbose" {
				viper.Set(setKey, setValue == "true" || setValue == "True")
			} else {
				viper.Set(setKey, setValue)
			}

			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error writing config:", err)
				return
			}

			fmt.Printf("Config updated: %s = %s\n", setKey, setValue)
			Log.AppendLog(fmt.Sprintf("Updated config: %s = %s", setKey, setValue))
			return
		}

		// Show all config if no --set
		fmt.Println("Current Configuration:")
		for _, key := range viper.AllKeys() {
			fmt.Printf("  %s: %v\n", key, viper.Get(key))
		}
	},
}

func init() {
	configCmd.Flags().StringVar(&setKey, "set-key", "", "Configuration key to set (e.g., treelimit)")
	configCmd.Flags().StringVar(&setValue, "set-value", "", "Value to set for the configuration key")

	rootCmd.AddCommand(configCmd)
}
