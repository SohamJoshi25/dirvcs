// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package services

import (
	Path "dirvcs/internal/data/path"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func AppendLog(entry string) error {
	logPath := Path.LOGS_PATH

	var logs []string

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	logs = append(logs, time.Now().String()+" : "+entry)

	newData, _ := json.MarshalIndent(logs, "", "  ")
	return os.WriteFile(logPath, newData, 0644)
}

func PrintLogs() {
	logPath := Path.LOGS_PATH

	var logs []string

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	fmt.Print("=== Logs ===\n")

	for _, log := range logs {
		fmt.Printf("%s\n", log)
	}
}
