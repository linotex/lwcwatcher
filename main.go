package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"lwcWatcher/src/config"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "0.1"

func main() {

	isHelp := flag.Bool("h", false, "Show this screen")
	first := flag.Bool("first", false, "Run watcher first time")
	file := flag.String("f", "", "Changed file")
	version := flag.Bool("v", false, "Version of watcher")

	flag.Parse()

	if *isHelp {
		flag.Usage()
		return
	}

	if *version {
		fmt.Println(VERSION)
		return
	}

	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	project := config.LoadConfig(currentDir)

	watchPackage := project.GetWatchPackage()
	defaultPackage := project.GetDefaultPackage()

	if watchPackage == "" {
		log.Fatal("Nothing to watch")
	}

	if defaultPackage == "" {
		log.Fatal("No default package")
	}

	if *first {
		CopyAllLwc(currentDir, watchPackage, defaultPackage)
		return
	}

	if *file == "" {
		return
	}

	if !IsWatchPackage(watchPackage, currentDir, *file) {
		return
	}

	CopyLwc(defaultPackage, currentDir, *file)
}

func IsWatchPackage(watchPackage string, projectDir string, file string) bool {
	lwcPath := projectDir + "/" + watchPackage + "/main/default/lwc"
	return strings.Index(file, lwcPath) == 0
}

func CopyLwc(defaultPackage string, projectDir string, file string) {

	fmt.Println("File: ", file)

	filePathParts := strings.Split(file, "/")

	fileName := filePathParts[len(filePathParts) - 1]
	componentName := filePathParts[len(filePathParts) - 2]

	targetDir := projectDir + "/" + defaultPackage + "/main/default/lwc/" + componentName
	copyFile(file, targetDir + "/" + fileName)
}

func CopyAllLwc(currentDir string, watchPackage string, defaultPackage string) {
	var componentPaths []string
	sourceDir := currentDir + "/" + watchPackage + "/main/default"
	targetDir := currentDir + "/" + defaultPackage + "/main/default/lwc"

	componentPaths = readLwcDir(sourceDir, componentPaths)
	fmt.Println("Count", len(componentPaths))

	for _, f := range componentPaths {
		arr := strings.Split(f, "/")
		dirName := arr[len(arr) - 1]
		copyDir(f, targetDir + "/" + dirName)
	}

	fmt.Println("Done.")
}

func readLwcDir(dir string, componentPaths []string) []string {
	lwcDir := dir + "/lwc"

	if !dirIsExist(lwcDir) {
		return append(componentPaths, dir)
	} else {
		folders := getListFiles(lwcDir, true)
		for _, path := range folders {
			componentPaths = readLwcDir(path, componentPaths)
		}
		return componentPaths
	}
}

func dirIsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getListFiles(path string, isDir bool) []string {
	var paths []string

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	list, err := f.Readdir(-1)
	_ = f.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range list {
		if isDir {
			if v.IsDir() {
				paths = append(paths, path + "/" + v.Name())
			}
		} else {
			paths = append(paths, path + "/" + v.Name())
		}
	}

	return paths
}

func copyDir(source string, target string)  {
	files := getListFiles(source, false)
	for _, f := range files {
		stat, _ := os.Stat(f)
		copyFile(f, target + "/" + stat.Name())
	}
}

func copyFile(source string, target string) {
	targetDir := filepath.Dir(target)
	_, err := os.Stat(target)

	if os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	fin, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	fout, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		log.Fatal(err)
	}
}