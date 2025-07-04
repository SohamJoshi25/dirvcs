// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package ignore

import (
	Path "dirvcs/internal/data/path"
	"fmt"
	"os"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

var Ignore *ignore.GitIgnore

func init() {
	ignore, err := ignore.CompileIgnoreFile(Path.IGNORE_PATH)
	if err != nil {
		Ignore = nil
	} else {
		Ignore = ignore
	}
}

func PrintIgnore() {
	content, err := os.ReadFile(Path.IGNORE_PATH)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logs := strings.Split(string(content), "\n")
	fmt.Printf("=== Ignore ===\n")

	for _, log := range logs {
		fmt.Println(log)
	}
}

func ApendIgnore(messages []string) {

	file, err := os.OpenFile(Path.IGNORE_PATH, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	content, err := os.ReadFile(Path.IGNORE_PATH)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(content) == 0 {
		_, err = file.WriteString(strings.Join(messages, "\n"))
	} else {
		_, err = file.WriteString("\n" + strings.Join(messages, "\n"))
	}

	// Write the text to the file

	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}

func RemoveIgnore(target string) {

	content, err := os.ReadFile(Path.IGNORE_PATH)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var lines []string
	logs := strings.Split(string(content), "\n")

	for _, log := range logs {
		if log != target {
			lines = append(lines, log)
		}
	}

	os.WriteFile(Path.IGNORE_PATH, []byte(strings.Join(lines, "\n")), 0644)
}
