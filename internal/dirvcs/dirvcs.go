package dirvcs

//Speed : 13.6 GB (14,66,08,95,020 bytes) 2,31,946 Files, 25,984 Folders in 15.1408093s

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
	"github.com/spf13/viper"

	Path "dirvcs/internal/data/path"
	Color "dirvcs/internal/services/color"
	Ignore "dirvcs/internal/services/ignore"
	Log "dirvcs/internal/services/logging"
	TLog "dirvcs/internal/services/treelogs"
	Struct "dirvcs/internal/structs"
)

var GINDENT string
var VERBOSE bool

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

	if VERBOSE {
		name = fmt.Sprintf("%s %s %s", name, node.ModificationTime, node.Hash)
	}

	fmt.Println(Color.Color(indent+name, color))

	for _, child := range node.Children {
		printTree(child, indent+GINDENT, color)
	}
}

func saveTreeJson(root *Struct.FileNode, path string) error {
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func CompareLevel(oldNode *Struct.FileNode, newNode *Struct.FileNode, indent string) {

	if oldNode.Hash == newNode.Hash {
		return
	}

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
				printTree(oldChild, indent, Color.Red)
			} else {

				if VERBOSE {
					fmt.Printf(Color.Red+"%s'%s' was Deleted. %s %s\n"+Color.Reset, indent, oldChild.Name, oldChild.ModificationTime, oldChild.Hash)
				} else {
					fmt.Printf(Color.Red+"%s'%s' was Deleted.\n"+Color.Reset, indent, oldChild.Name)
				}

			}
		} else if status == 1 {

			if newChild.IsDir {

				if VERBOSE {
					fmt.Printf(Color.Gray+"%s'%s' %s %s\n"+Color.Reset, indent, oldChild.Name, oldChild.ModificationTime, oldChild.Hash)
				} else {
					fmt.Printf(Color.Gray+"%s'%s'\n"+Color.Reset, indent, newChild.Name)
				}
			} else {
				if VERBOSE {
					fmt.Printf(Color.Yellow+"%s'%s' has SOME CHANGES. %s %s\n"+Color.Reset, indent, oldChild.Name, oldChild.ModificationTime, oldChild.Hash)
				} else {
					fmt.Printf(Color.Yellow+"%s'%s' has SOME CHANGES.\n"+Color.Reset, indent, newChild.Name)
				}

			}

			childNodeExist[childIndex] = true

			if oldChild.IsDir {
				CompareLevel(oldChild, newChild, indent+GINDENT)
			}

		} else {
			childNodeExist[childIndex] = true
		}
	}

	for index, isCounted := range childNodeExist {
		if !isCounted {

			if newNode.Children[index].IsDir {
				printTree(newNode.Children[index], indent, Color.Green)
			} else {

				if VERBOSE {
					fmt.Printf(Color.Green+"%s'%s' was Created. %s %s\n"+Color.Reset, indent, newNode.Name, newNode.ModificationTime, newNode.Hash)
				} else {
					fmt.Printf(Color.Green+"%s'%s' was Created.\n"+Color.Reset, indent, newNode.Children[index].Name)
				}
			}
		}

	}
}

//Public

func GenerateTree(BASE_PATH, message string) {
	initConfig()

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
	treePath, errabs := filepath.Abs(path.Join(Path.TREES_PATH, fmt.Sprintf(`%s.gz`, uuid)))
	if errabs != nil {
		log.Fatalln("Unable to get Absolute Path")
	}

	DirRecursveInfo(rootNode)
	SaveTree(rootNode, treePath)

	var TreeLog *Struct.TreeLog = &Struct.TreeLog{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		TreePath:  treePath,
		TreeHash:  rootNode.Hash,
		TreeId:    uuid.String(),
	}

	TLog.AppendLog(TreeLog)
	Log.AppendLog(fmt.Sprintf("tree generated %s %s", uuid, message))

	elapsed := time.Since(start)

	fmt.Printf("\nDirectory Persisted '%s' %s", message, uuid)
	fmt.Printf("\nTime took %s", elapsed)

	TLog.LimitTree()

}

func PrintTree(index int) {
	initConfig()

	if index == 0 {

		BASE_PATH := "."

		var newTree *Struct.FileNode

		info, err := os.Stat(BASE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		absPath, err := filepath.Abs(BASE_PATH)

		if err != nil {
			log.Fatal("Absolute File Parsing Error")
		}

		newTree = &Struct.FileNode{
			Name:             absPath,
			Path:             absPath,
			Depth:            0,
			IsDir:            true,
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), BASE_PATH, info.ModTime().Format(time.RFC3339))),
			Children:         []*Struct.FileNode{},
		}

		DirRecursveInfo(newTree)
		printTree(newTree, "", Color.Gray)

		return
	} else {
		index--
	}

	treelog, err := TLog.LastLogIdx(index)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	TLog.PrintTreeLog(treelog)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootNode, err := LoadTree(treelog.TreePath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	start := time.Now()

	printTree(rootNode, "", Color.Gray)

	elapsed := time.Since(start)
	fmt.Printf("\nTime took %s", elapsed)
}

func PrintTreeUUID(uuid string) {
	initConfig()

	treelog, err := TLog.GetByUuid(uuid)

	TLog.PrintTreeLog(treelog)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rootNode, err := LoadTree(treelog.TreePath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	start := time.Now()

	printTree(rootNode, "", Color.Gray)

	elapsed := time.Since(start)
	fmt.Printf("\nTime took %s", elapsed)
}

func ExportTree(uuid string, filepath string) {
	treelog, err := TLog.GetByUuid(uuid)

	TLog.PrintTreeLog(treelog)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rootNode, err := LoadTree(treelog.TreePath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	start := time.Now()

	saveTreeJson(rootNode, filepath)

	elapsed := time.Since(start)
	fmt.Printf("\nTime took %s", elapsed)
}

func CompareTree(oldId, newId string) {
	initConfig()

	var oldTreeLog *Struct.TreeLog
	var err1 error

	if oldId == "" {
		oldTreeLog, err1 = TLog.LastLogIdx(0)
	} else {
		oldTreeLog, err1 = TLog.GetByUuid(oldId)
	}

	if err1 != nil {
		fmt.Printf("%v\n", err1)
		os.Exit(1)
	}

	var newTree *Struct.FileNode

	if newId == "" {

		BASE_PATH := "."

		info, err := os.Stat(BASE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		absPath, err := filepath.Abs(BASE_PATH)

		if err != nil {
			log.Fatal("Absolute File Parsing Error")
		}

		newTree = &Struct.FileNode{
			Name:             absPath,
			Path:             absPath,
			Depth:            0,
			IsDir:            true,
			ModificationTime: info.ModTime().Format(time.RFC3339),
			Size:             uint64(info.Size()),
			Hash:             SHA256(fmt.Sprintf("%s %s %s", info.Name(), BASE_PATH, info.ModTime().Format(time.RFC3339))),
			Children:         []*Struct.FileNode{},
		}

		DirRecursveInfo(newTree)

	} else {
		newTreeLog, err1 := TLog.GetByUuid(newId)

		if err1 != nil {
			fmt.Printf("UUID not found %s", newId)
		}

		newTree, err1 = LoadTree(newTreeLog.TreePath)

		if err1 != nil {
			fmt.Printf("Error Loading Tree %v", err1)
		}
	}

	oldTree, err2 := LoadTree(oldTreeLog.TreePath)

	if err1 != nil || err2 != nil {
		fmt.Printf("Could Not Load Tree ")
	}

	if oldTree.Hash != newTree.Hash {
		fmt.Printf(Color.Gray+"'%s'\n"+Color.Reset, newTree.Path)
		CompareLevel(oldTree, newTree, GINDENT)
	} else {
		fmt.Println("No Changes Found")
	}

}

func CompareTreePath(oldPath, newPath string) {
	initConfig()

	newTree, err1 := LoadTree(newPath)
	oldTree, err2 := LoadTree(oldPath)

	if err1 != nil || err2 != nil {
		fmt.Printf("Could Not Load Tree ")
		if err1 != nil {
			fmt.Println(err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		os.Exit(1)
	}

	if oldTree.Hash != newTree.Hash {
		fmt.Printf(Color.Gray+"'%s'\n"+Color.Reset, newTree.Path)
		CompareLevel(oldTree, newTree, GINDENT)
	} else {
		fmt.Println("No Changes Found")
	}
}

func initConfig() {
	GINDENT = viper.GetString("indent")
	if GINDENT == "" {
		GINDENT = "|--"
	}
	VERBOSE = viper.GetBool("verbose")
}
