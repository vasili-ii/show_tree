package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var excludeNames map[string]bool = map[string]bool{
	".DS_Store":   true,
	".vscode":     true,
	"launch.json": true,
}

func getNameWithSize(baseName string, info os.FileInfo) (string, error) {
	if info.IsDir() {
		return baseName, nil
	}

	var fileSize string
	if info.Size() > 0 {
		fileSize = fmt.Sprintf(" (%vb)", info.Size())
	} else {
		fileSize = " (empty)"
	}

	return baseName + fileSize, nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		baseName := filepath.Base(path)
		if _, inExcluded := excludeNames[baseName]; inExcluded {
			return nil
		}

		pathArr := strings.Split(path, string(os.PathSeparator))
		pathLength := len(pathArr)
		var prefix string = "├───"
		if pathLength > 0 {
			// fmt.Printf("pathLength: %v, path: %v\n", pathLength, path)
			prefix = strings.Repeat("│  ", (pathLength - 1))
			prefix += "├───"
		}

		pointName, err := getNameWithSize(baseName, info)
		if err == nil {
			fmt.Printf("%v%v\n", prefix, pointName)
		}
		return nil
	})

	return err
}

func main() {
	fmt.Println("\n\n\n")
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
