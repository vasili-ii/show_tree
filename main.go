package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var excludeNames map[string]bool = map[string]bool{
	".DS_Store":   true,
	".vscode":     true,
	"launch.json": true,
	".git":        true,
	"__debug_bin": true,
}

func getNameWithSize(info os.FileInfo) (string, error) {
	if info.IsDir() {
		return info.Name(), nil
	}

	var fileSize string
	if info.Size() > 0 {
		fileSize = fmt.Sprintf(" (%vb)", info.Size())
	} else {
		fileSize = " (empty)"
	}

	return info.Name() + fileSize, nil
}

func dirTreeSecond(out io.Writer, path string, printFiles bool) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Prevent panic by handling failure accessing a path %q: %v\n", path, err)
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

		pointName, err := getNameWithSize(info)
		if err == nil {
			fmt.Fprintf(out, "%v%v\n", prefix, pointName)
		}
		return nil
	})

	return err
}

func dirTreeRecursive(out io.Writer, path string, printFiles bool, prefix string) error {

	curprefix := prefix

	dirContent, _ := ioutil.ReadDir(path)
	// fmt.Fprint(out, "├───")

	for idx := range dirContent {
		flInfo := dirContent[idx]
		prefix = curprefix + "├───"

		if idx == len(dirContent) {
			// по последнему файлу в директории ставим такой значок:
			prefix = curprefix + "└───"
		}
		if _, inExcluded := excludeNames[flInfo.Name()]; inExcluded {
			continue
		}
		nm, _ := getNameWithSize(flInfo)
		fmt.Fprintf(out, "%v%v\n", prefix, nm)

		if flInfo.IsDir() {
			prefixToRecursive := curprefix + "|\t"
			dirTreeRecursive(out, flInfo.Name(), printFiles, prefixToRecursive)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	dirTreeRecursive(out, path, printFiles, "")
	return nil
}

func main() {
	out := os.Stdout
	var path string
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		// panic("usage go run main.go . [-f]")
		path = "."
	} else {
		path = os.Args[1]
	}
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
