package dirvcs

//Speed : 13.6 GB (14,66,08,95,020 bytes) 2,31,946 Files, 25,984 Folders in 7.8302912s

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	Path "dirvcs/internal/data/path"
	Color "dirvcs/internal/services/color"
	Ignore "dirvcs/internal/services/ignore"
	"dirvcs/internal/services/treelogs"
	Struct "dirvcs/internal/structs"
)

func SHA256(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashSum := hasher.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func DirRecursveInfo(node *Struct.FileNode) {
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		log.Println("Skipping:", node.Path, "due to error:", err)
		return
	}

	var hash string = node.Hash
	var size uint64 = node.Size

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			log.Println("Error getting info for", e.Name(), ":", err)
			continue
		}
		path := filepath.Join(node.Path, info.Name())

		if Ignore.Ignore != nil && Ignore.Ignore.MatchesPath(path) {
			continue
		}

		var children []*Struct.FileNode
		if info.IsDir() {
			children = []*Struct.FileNode{}
		} else {
			children = nil
		}

		childnode := &Struct.FileNode{
			Name:             info.Name(),
			Path:             path,
			Depth:            node.Depth + 1,
			IsDir:            info.IsDir(),
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), path, info.ModTime().Format(time.RFC3339))),
			Children:         children,
		}

		if info.IsDir() {
			DirRecursveInfo(childnode)
		}

		size += childnode.Size
		hash = SHA256(fmt.Sprintf("%s %s", hash, childnode.Hash))

		node.Children = append(node.Children, childnode) // Append Child Node Refrence to Parent node children array

	}

	node.Hash = hash
	node.Size = size

}

func SaveTree(root *Struct.FileNode, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	encoder := json.NewEncoder(gzWriter)
	return encoder.Encode(root)
}

func LoadTree(path string) (*Struct.FileNode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	var root Struct.FileNode
	decoder := json.NewDecoder(gzReader)
	err = decoder.Decode(&root)
	return &root, err
}

func printTree(node *Struct.FileNode, indent string, color string) {
	var name string

	switch color {
	case Color.Red:
		name = fmt.Sprintf("'%s' was DELETED.", node.Name)
	case Color.Yellow:
		name = fmt.Sprintf("'%s' has SOME CHANGES.", node.Name)
	case Color.Green:
		name = fmt.Sprintf("'%s' was CREATED.", node.Name)
	default:
		name = fmt.Sprintf("'%s'", node.Name)
	}
	fmt.Println(color + indent + name + Color.Reset)
	for _, child := range node.Children {
		printTree(child, indent+"|---", color)
	}
}

func CompareLevel(oldNode *Struct.FileNode, newNode *Struct.FileNode, indent string) {
	//Assumption : There are changes in old and new node 's Hash (Some Changes Exist)

	getChildFromNewNode := func(hash string, filename string) (byte, *Struct.FileNode, int) {
		// 0 -> No Node Found
		// 1 -> Node Found with same File Name but mismatching Hash
		// 2 -> Same Node found with matching hash

		for idx, newChild := range newNode.Children {
			if newChild.Name == filename {
				if newChild.Hash == hash {
					return 2, newChild, idx
				} else {
					return 1, newChild, idx
				}
			}
		}

		return 0, nil, -1
	}

	childNodeExist := make([]bool, len(newNode.Children))

	for _, oldChild := range oldNode.Children {
		status, newChild, childIndex := getChildFromNewNode(oldChild.Hash, oldChild.Name)

		if status == 0 {
			//fmt.Printf(Green+"%s'%s' was CREATED.\n"+Reset, indent, oldChild.Name)
			if oldChild.IsDir {
				printTree(oldChild, indent, Color.Green)
			} else {
				fmt.Printf(Color.Green+"%s'%s' was CREATED.\n"+Color.Reset, indent, oldChild.Name)
			}
		} else if status == 1 {

			if newChild.IsDir {
				fmt.Printf(Color.Gray+"%s'%s'\n"+Color.Reset, indent, newChild.Name)
			} else {
				fmt.Printf(Color.Yellow+"%s'%s' has SOME CHANGES.\n"+Color.Reset, indent, newChild.Name)
			}

			childNodeExist[childIndex] = true

			if oldChild.IsDir {
				CompareLevel(oldChild, newChild, indent+"|---")
			}

		} else {
			childNodeExist[childIndex] = true
		}
	}

	for index, isCounted := range childNodeExist {
		if !isCounted {

			if newNode.Children[index].IsDir {
				printTree(newNode.Children[index], indent, Color.Red)
			} else {
				fmt.Printf(Color.Red+"%s'%s' was DELETED.\n"+Color.Reset, indent, newNode.Children[index].Name)
			}
		}

	}
}

//Public

func GenerateTree(BASE_PATH, message string) {

	info, err := os.Stat(BASE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal("Input Directory Path cannot be a file.")
	}

	absPath, err := filepath.Abs(BASE_PATH)

	if err != nil {
		log.Fatal("Absolute File Parsing Error")
	}

	uuid := uuid.New()

	var rootNode *Struct.FileNode = &Struct.FileNode{
		Name:             absPath,
		Path:             absPath,
		Depth:            0,
		IsDir:            true,
		ModificationTime: info.ModTime().Format(time.RFC3339),
		Size:             uint64(info.Size()),
		Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), BASE_PATH, info.ModTime().Format(time.RFC3339))),
		Children:         []*Struct.FileNode{},
	}

	start := time.Now()

	DirRecursveInfo(rootNode)
	SaveTree(rootNode, path.Join(Path.TREES_PATH, fmt.Sprintf(`%s.gz`, uuid)))

	var TreeLog *Struct.TreeLog = &Struct.TreeLog{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		TreePath:  path.Join(Path.TREES_PATH, fmt.Sprintf(`%s.gz`, uuid)),
		TreeHash:  rootNode.Hash,
	}

	treelogs.AppendLog(TreeLog)

	elapsed := time.Since(start)

	fmt.Printf("\nDirectory Persisted '%s' %s", message, uuid)
	fmt.Printf("\nTime took %s", elapsed)

}

func PrintTree(index int) {

	treelog, err := treelogs.LastLog(index)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	rootNode, err := LoadTree(treelog.TreePath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	start := time.Now()

	printTree(rootNode, "  ", Color.Gray)

	elapsed := time.Since(start)
	fmt.Printf("\nTime took %s", elapsed)
}

func main() {
	fmt.Println("===== DIRVCS =====\n\nEnter Your Choice")
	fmt.Println("1. Generate Tree")
	fmt.Println("2. Display Tree")
	fmt.Println("3. Compare Tree")

	var choice int
	fmt.Scan(&choice)

	if choice == 1 {

	} else if choice == 2 {

	} else if choice == 3 {

		var OLD_PATH string
		var NEW_PATH string

		fmt.Printf("\nEnter Path of generated tree : ")
		fmt.Scan(&OLD_PATH)

		OLD_PATH = filepath.Join(OLD_PATH, "tree.json")

		fmt.Printf("\nEnter Path of New Directory : ")
		fmt.Scan(&NEW_PATH)

		fmt.Printf("\nDirectory Change Log \n\n")

		info, err := os.Stat(".")
		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			log.Fatal("Input Directory Path cannot be a file.")
		}

		var oldTree *Struct.FileNode = &Struct.FileNode{
			Name:             info.Name(),
			Path:             NEW_PATH,
			Depth:            0,
			IsDir:            true,
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), NEW_PATH, info.ModTime().Format(time.RFC3339))),
			Children:         []*Struct.FileNode{},
		}

		start := time.Now()

		newTree, err := LoadTree(OLD_PATH)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		DirRecursveInfo(oldTree)

		CompareLevel(oldTree, newTree, "")

		elapsed := time.Since(start)
		fmt.Printf("\nTime took %s", elapsed)

	} else {
		fmt.Println("Wrong Choice")
	}

	os.Exit(0)
}
