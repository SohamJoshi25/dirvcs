package treelogs

import (
	Path "dirvcs/internal/data/path"
	Struct "dirvcs/internal/structs"
	"encoding/json"
	"fmt"
	"os"
)

func AppendLog(root *Struct.TreeLog) error {
	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	logs = append(logs, root)

	newData, _ := json.MarshalIndent(logs, "", "  ")
	return os.WriteFile(logPath, newData, 0644)
}

func LastLog(index int) (*Struct.TreeLog, error) {

	if index < 0 {
		return nil, fmt.Errorf("Index cannot be less than 0.")
	}

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if index >= len(logs) {
		return nil, fmt.Errorf(`Index cannot be greater than %d`, len(logs))
	}

	return logs[len(logs)-1-index], nil
}

func DeleteLog(index int) error {

	if index < 0 {
		return fmt.Errorf("Index cannot be less than 0.")
	}

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if index >= len(logs) {
		return fmt.Errorf(`Index cannot be greater than %d`, len(logs))
	}

	logs = append(logs[:index], logs[index+1:]...)
	newData, _ := json.MarshalIndent(logs, "", "  ")

	return os.WriteFile(logPath, newData, 0644)
}
