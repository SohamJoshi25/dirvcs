// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package treelogs

import (
	Path "dirvcs/internal/data/path"
	Logs "dirvcs/internal/services/logging"
	Struct "dirvcs/internal/structs"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
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

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("please persist atleast once to compare versions")
	}

	if index < 0 {
		return nil, fmt.Errorf("index cannot be less than 0")
	}

	if index >= len(logs) {
		return nil, fmt.Errorf(`index cannot be greater than %d`, len(logs))
	}

	return logs[len(logs)-1-index], nil
}

func GetByUuid(uuid string) (*Struct.TreeLog, error) {

	logPath := Path.TREE_LOG_PATH
	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("please persist atleast once to export version")
	}

	if uuid == "" {
		return logs[len(logs)-1], nil
	}

	for idx := range logs {
		if logs[idx].TreeId == uuid {
			return logs[idx], nil
		}
	}

	return nil, fmt.Errorf("tree not found")
}

func DeleteLogIdx(index int) error {

	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if len(logs) == 0 {
		return fmt.Errorf("please persist atleast once to delete tree")
	}

	if index < 0 {
		return fmt.Errorf("index cannot be less than 0")
	}

	if index >= len(logs) {
		return fmt.Errorf(`index cannot be greater than %d`, len(logs))
	}

	er := os.Remove(logs[index].TreePath)
	if er != nil {
		fmt.Println(er)
		os.Exit(1)
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

	er := os.Remove(logs[index].TreePath)
	if er != nil {
		fmt.Println(er)
		os.Exit(1)
	}

	logs = append(logs[:index], logs[index+1:]...)
	newData, _ := json.MarshalIndent(logs, "", "  ")

	err := os.WriteFile(logPath, newData, 0644)
	return err
}

func LimitTree() {
	logPath := Path.TREE_LOG_PATH

	var logs []*Struct.TreeLog

	if data, err := os.ReadFile(logPath); err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file not found")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	TreeLimit := viper.GetInt("treelimit")

	if len(logs) > TreeLimit {
		DeleteLogUuid(logs[0].TreeId)
		Logs.AppendLog(fmt.Sprintf("Trees %s Pruned to the set tree limit %d", logs[0].TreeId, TreeLimit))
	}
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
