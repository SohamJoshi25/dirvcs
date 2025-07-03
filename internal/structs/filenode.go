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

func ToFileNodeChanges(filenode *FileNode, Operation string) *FileNodeChanges {
	FileNodeChanges := &FileNodeChanges{
		Operation:        Operation,
		Name:             filenode.Name,
		Path:             filenode.Path,
		Depth:            filenode.Depth,
		IsDir:            filenode.IsDir,
		ModificationTime: filenode.ModificationTime,
		Size:             filenode.Size,
		Hash:             filenode.Hash,
		Children:         []*FileNodeChanges{},
	}
	return FileNodeChanges
}

type FileNodeChanges struct {
	Operation        string             `json:"operation"`
	Name             string             `json:"name"`
	Path             string             `json:"path"`
	Depth            int                `json:"depth"`
	IsDir            bool               `json:"isDir"`
	ModificationTime string             `json:"modificationTime"`
	Size             uint64             `json:"size"`
	Hash             string             `json:"hash"`
	Children         []*FileNodeChanges `json:"children"`
}
