package structs

type TreeLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	TreePath  string `json:"treepath"`
	TreeHash  string `json:"treehash"`
}
