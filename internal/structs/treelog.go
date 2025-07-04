// Licensed under the Polyform Noncommercial License 1.0.0
// You may use, copy, modify, and distribute this software for noncommercial purposes.
// See LICENSE for details

package structs

type TreeLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	TreePath  string `json:"treepath"`
	TreeHash  string `json:"treehash"`
	TreeId    string `json:"treeid"`
}
