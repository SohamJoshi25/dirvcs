package structs

type TreeLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	TreeName  string `json:"treename"`
	TreeHash  string `json:"treehash"`
}
