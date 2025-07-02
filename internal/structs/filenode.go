package structs

type FileNode struct {
	Name             string      `json:"name"`
	Path             string      `json:"path"`
	Depth            int         `json:"depth"`
	IsDir            bool        `json:"isDir"`
	ModificationTime string      `json:"modificationTime"`
	Size             uint64      `json:"size"`
	Hash             string      `json:"hash"`
	Children         []*FileNode `json:"children"`
}
