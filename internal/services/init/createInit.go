// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package services

import (
	path "dirvcs/internal/data/path"
	Log "dirvcs/internal/services/logging"

	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func CreateInit() {

	_, err := os.Stat(path.BASE_PATH)
	if err == nil {
		fmt.Println("dirvcs is already initialised")
		os.Exit(0)
	}

	if err := os.MkdirAll(path.TREES_PATH, os.ModePerm); err != nil {
		fmt.Printf("unable to create trees directory: %v", err)
		os.Exit(1)
	}
	if err := os.WriteFile(path.TREE_LOG_PATH, []byte("[]"), 0755); err != nil {
		fmt.Printf("unable to create tree log file: %v", err)
		os.Exit(1)
	}

	if err := os.WriteFile(path.IGNORE_PATH, []byte(".*\n"), 0755); err != nil {
		fmt.Printf("unable to create ignore file: %v", err)
		os.Exit(1)
	}

	if err := viper.SafeWriteConfig(); err != nil {
		log.Fatalf("Failed to create config file: %v", err)
		os.Exit(1)
	}

	if err := os.WriteFile(path.LOGS_PATH, []byte("[]"), 0755); err != nil {
		fmt.Printf("unable to create log file: %v", err)
		os.Exit(1)
	}

	Log.AppendLog("dirvcs was initialised.")

	fmt.Println("dirvcs initialised")
}

func CheckInit() {
	_, err := os.Stat(path.BASE_PATH)
	if err != nil {
		fmt.Println("dirvcs not initialised in this directory")
		os.Exit(1)
	}
}
