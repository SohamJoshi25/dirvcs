package cmd

import (
	"fmt"
	"log"
	"strconv"

	Init "dirvcs/internal/services/init"
	Log "dirvcs/internal/services/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setKey, setValue string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update configuration",
	Long:  "View all current configuration values, or update a specific config key like tree_limit.",
	Run: func(cmd *cobra.Command, args []string) {

		Init.CheckInit()

		// Update if --set is used
		if setKey != "" && setValue != "" {
			switch setKey {
			case "treelimit":
				val, err := strconv.Atoi(setValue)
				if err != nil {
					fmt.Println("treelimit must be an integer")
					return
				}
				viper.Set("treelimit", val)
			default:
				viper.Set(setKey, setValue)
			}

			if err := viper.WriteConfig(); err != nil {
				log.Fatalf("Failed to write config: %v", err)
			}
			fmt.Printf("Updated config: %s = %s\n", setKey, setValue)
			Log.AppendLog(fmt.Sprintf("Updated config: %s = %s", setKey, setValue))
			return
		}

		// If no --set, display all config values
		fmt.Println("Current Configuration:")
		for _, key := range viper.AllKeys() {
			fmt.Printf("  %s: %v\n", key, viper.Get(key))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVar(&setKey, "set-key", "", "Configuration key to set (e.g., treelimit)")
	configCmd.Flags().StringVar(&setValue, "set-value", "", "Value to set for the configuration key")
}
