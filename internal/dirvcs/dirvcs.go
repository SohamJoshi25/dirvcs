package dirvcs

//Speed : 13.6 GB (14,66,08,95,020 bytes) 2,31,946 Files, 25,984 Folders in 7.8302912s

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ignore "github.com/sabhiram/go-gitignore"
)

var ign *ignore.GitIgnore

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

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

func SHA256(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashSum := hasher.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func dirRecursveInfo(node *FileNode) {
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

		if ign != nil && ign.MatchesPath(path) {
			continue
		}
		// fmt.Println(node.Path)
		// fmt.Println(node.Depth)
		// fmt.Println(info.Name())
		// fmt.Println(info.IsDir())
		// fmt.Println(info.ModTime())
		// fmt.Println(info.Mode())
		// fmt.Println(info.Size())
		// fmt.Println(info.Sys())
		// fmt.Println()

		var children []*FileNode
		if info.IsDir() {
			children = []*FileNode{}
		} else {
			children = nil
		}

		childnode := &FileNode{
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
			dirRecursveInfo(childnode)
		}

		size += childnode.Size
		hash = SHA256(fmt.Sprintf("%s %s", hash, childnode.Hash))

		node.Children = append(node.Children, childnode) // Append Child Node Refrence to Parent node children array

	}

	node.Hash = hash
	node.Size = size

}

func saveTree(root *FileNode, path string) error {
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func printTree(node *FileNode, indent string, color string) {
	var name string

	switch color {
	case Red:
		name = fmt.Sprintf("'%s' was DELETED.", node.Name)
	case Yellow:
		name = fmt.Sprintf("'%s' has SOME CHANGES.", node.Name)
	case Green:
		name = fmt.Sprintf("'%s' was CREATED.", node.Name)
	default:
		name = fmt.Sprintf("'%s'", node.Name)
	}
	fmt.Println(color + indent + name + Reset)
	for _, child := range node.Children {
		printTree(child, indent+"|---", color)
	}
}

func loadTree(path string) (*FileNode, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root FileNode
	err = json.Unmarshal(data, &root)
	return &root, err
}

func compareLevel(oldNode *FileNode, newNode *FileNode, indent string) {
	//Assumption : There are changes in old and new node 's Hash (Some Changes Exist)

	getChildFromNewNode := func(hash string, filename string) (byte, *FileNode, int) {
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
				printTree(oldChild, indent, Green)
			} else {
				fmt.Printf(Green+"%s'%s' was CREATED.\n"+Reset, indent, oldChild.Name)
			}
		} else if status == 1 {

			if newChild.IsDir {
				fmt.Printf(Gray+"%s'%s'\n"+Reset, indent, newChild.Name)
			} else {
				fmt.Printf(Yellow+"%s'%s' has SOME CHANGES.\n"+Reset, indent, newChild.Name)
			}

			childNodeExist[childIndex] = true

			if oldChild.IsDir {
				compareLevel(oldChild, newChild, indent+"|---")
			}

		} else {
			//fmt.Printf("%s'%s' has NO CHANGES.\n", indent, newChild.Name)
			childNodeExist[childIndex] = true
		}
	}

	for index, isCounted := range childNodeExist {
		if !isCounted {

			if newNode.Children[index].IsDir {
				printTree(newNode.Children[index], indent, Red)
			} else {
				fmt.Printf(Red+"%s'%s' was DELETED.\n"+Reset, indent, newNode.Children[index].Name)
			}
		}

	}
}

func main() {
	fmt.Println("===== DIRVCS =====\n\nEnter Your Choice")
	fmt.Println("1. Generate Tree")
	fmt.Println("2. Display Tree")
	fmt.Println("3. Compare Tree")

	ignore, err := ignore.CompileIgnoreFile(".dirignore")
	if err != nil {
		ign = nil
	} else {
		ign = ignore
	}

	var choice int
	fmt.Scan(&choice)

	if choice == 1 {

		var BASE_PATH string
		var OUTPUT_PATH string

		fmt.Printf("\nEnter Base Path to generate tree : ")
		fmt.Scan(&BASE_PATH)

		fmt.Printf("\nEnter Output Path of Generated Tree : ")
		fmt.Scan(&OUTPUT_PATH)

		OUTPUT_PATH = filepath.Join(OUTPUT_PATH, "tree.json")

		info, err := os.Stat(BASE_PATH)
		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			log.Fatal("Input Directory Path cannot be a file.")
		}

		var rootNode *FileNode = &FileNode{
			Name:             info.Name(),
			Path:             BASE_PATH,
			Depth:            0,
			IsDir:            true,
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), BASE_PATH, info.ModTime().Format(time.RFC3339))),
			Children:         []*FileNode{},
		}

		start := time.Now()

		dirRecursveInfo(rootNode)
		saveTree(rootNode, OUTPUT_PATH)

		elapsed := time.Since(start)
		fmt.Printf("Time took %s", elapsed)

	} else if choice == 2 {

		fmt.Printf("\nEnter Path where tree.json exist to Print Tree : ")
		var OUTPUT_PATH string
		fmt.Scan(&OUTPUT_PATH)

		OUTPUT_PATH = filepath.Join(OUTPUT_PATH, "tree.json")
		rootNode, err := loadTree(OUTPUT_PATH)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		start := time.Now()

		printTree(rootNode, "  ", Gray)

		elapsed := time.Since(start)
		fmt.Printf("Time took %s", elapsed)

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

		var oldTree *FileNode = &FileNode{
			Name:             info.Name(),
			Path:             NEW_PATH,
			Depth:            0,
			IsDir:            true,
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), NEW_PATH, info.ModTime().Format(time.RFC3339))),
			Children:         []*FileNode{},
		}

		start := time.Now()

		newTree, err := loadTree(OLD_PATH)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		dirRecursveInfo(oldTree)

		compareLevel(oldTree, newTree, "")

		elapsed := time.Since(start)
		fmt.Printf("\nTime took %s", elapsed)

	} else {
		fmt.Println("Wrong Choice")
	}

	os.Exit(0)
}
