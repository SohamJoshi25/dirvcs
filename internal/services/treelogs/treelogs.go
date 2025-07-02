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

func LastLogIdx(index int) (*Struct.TreeLog, error) {

	if index < 0 {
		return nil, fmt.Errorf("index cannot be less than 0.")
	}

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if index >= len(logs) {
		return nil, fmt.Errorf(`index cannot be greater than %d`, len(logs)-1)
	}

	return logs[len(logs)-1-index], nil
}

func LastLogUuid(uuid string) (*Struct.TreeLog, error) {

	logPath := Path.TREE_LOG_PATH
	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	for idx := range logs {
		if logs[idx].TreeId == uuid {
			return logs[idx], nil
		}
	}

	return nil, fmt.Errorf("tree not found")
}

func DeleteLogIdx(index int) error {

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

func DeleteLogUuid(uuid string) error {

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	var index int = -1

	for idx := range logs {
		if logs[idx].TreeId == uuid {
			index = idx
			break
		}
	}

	if index == -1 {
		fmt.Printf("UUID Not Found")
		os.Exit(1)
	}

	logs = append(logs[:index], logs[index+1:]...)
	newData, _ := json.MarshalIndent(logs, "", "  ")

	return os.WriteFile(logPath, newData, 0644)
}

func PrintTreeLogs() {
	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	for idx, log := range logs {
		fmt.Printf("\n\n=== %d ===\n", idx)
		PrintTreeLog(log)
	}

}

func PrintTreeLog(treelog *Struct.TreeLog) {
	fmt.Printf("TreePath : %s\n", treelog.TreePath)
	fmt.Printf("TreeHash : %s\n", treelog.TreeHash)
	fmt.Printf("TimeStamp : %s\n", treelog.Timestamp)
	fmt.Printf("Persist Message : %s\n", treelog.Message)
	fmt.Printf("Tree UUID : %s\n", treelog.TreeId)
}
