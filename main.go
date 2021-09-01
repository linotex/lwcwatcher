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

const VERSION = "0.2"

var ProjectDir string
var WatchPackage string
var DefaultPackage string

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

	ProjectDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	project := config.LoadConfig(ProjectDir)

	WatchPackage = project.GetWatchPackage()
	DefaultPackage = project.GetDefaultPackage()

	if WatchPackage == "" {
		log.Fatal("Nothing to watch")
	}

	if DefaultPackage == "" {
		log.Fatal("No default package")
	}

	if *first {
		CopyAllLwc()
		return
	}

	if *file == "" {
		return
	}

	if !IsWatchPackage(*file) {
		return
	}

	fmt.Println("File: ", *file)

	if isLwcFile(*file) {
		CopyLwc(*file)
	} else if isStaticResourceFile(*file) {
		CopyStaticResource(*file)
	}
}

func watchPackagePath() string {
	return ProjectDir + "/" + WatchPackage + "/main/default"
}

func defaultPackagePath() string {
	return ProjectDir + "/" + DefaultPackage + "/main/default"
}

func IsWatchPackage(file string) bool {
	return strings.Index(file, watchPackagePath()) == 0
}

func isLwcFile(file string) bool {
	return strings.Index(file, watchPackagePath()+"/lwc") == 0
}

func isStaticResourceFile(file string) bool {
	return strings.Index(file, watchPackagePath()+"/staticresources") == 0
}

func CopyStaticResource(file string) {
	staticResourcesFile := strings.ReplaceAll(file, watchPackagePath(), "")

	targetFile := defaultPackagePath() + staticResourcesFile
	copyFile(file, targetFile)
}

func CopyLwc(file string) {
	filePathParts := strings.Split(file, "/")

	fileName := filePathParts[len(filePathParts)-1]
	componentName := filePathParts[len(filePathParts)-2]

	targetDir := defaultPackagePath() + "/lwc/" + componentName
	copyFile(file, targetDir+"/"+fileName)
}

func CopyAllLwc() {
	var componentPaths []string
	sourceDir := watchPackagePath()
	targetDir := defaultPackagePath() + "/lwc"

	componentPaths = readLwcDir(sourceDir, componentPaths)
	fmt.Println("Count", len(componentPaths))

	for _, f := range componentPaths {
		arr := strings.Split(f, "/")
		dirName := arr[len(arr)-1]
		copyDir(f, targetDir+"/"+dirName)
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
				paths = append(paths, path+"/"+v.Name())
			}
		} else {
			paths = append(paths, path+"/"+v.Name())
		}
	}

	return paths
}

func copyDir(source string, target string) {
	files := getListFiles(source, false)
	for _, f := range files {
		stat, _ := os.Stat(f)
		copyFile(f, target+"/"+stat.Name())
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
