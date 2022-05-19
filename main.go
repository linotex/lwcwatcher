package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"lwcWatcher/src/config"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const VERSION = "0.5"

var (
	ProjectDir     string
	WatchPackage   string
	DefaultPackage string
)

func main() {

	isHelp := flag.Bool("h", false, "Show this screen")
	first := flag.Bool("first", false, "Run watcher first time")
	file := flag.String("f", "", "Changed file")
	projectDir := flag.String("project-dir", "", "Project dir root")
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

	if *projectDir != "" {
		ProjectDir = *projectDir
	} else {
		ProjectDir, _ = filepath.Abs(path.Dir(os.Args[0]))
	}

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
		CopyAllStaticResources()
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

	if !isEqual(file, targetFile) {
		copyFile(file, targetFile)
	}
}

func CopyLwc(file string) {
	fileName := path.Base(file)
	componentName, componentDir := getComponentNameWithPath(file)

	relativePath, _ := filepath.Rel(componentDir, path.Dir(file))

	targetFile := defaultPackagePath() + "/lwc/" + componentName + "/" + relativePath + "/" + fileName

	if !isEqual(file, targetFile) {
		copyFile(file, targetFile)
	}
}

func CopyAllLwc() {
	var componentPaths []string
	sourceDir := watchPackagePath()
	targetDir := defaultPackagePath() + "/lwc"

	fmt.Println("Start copy LWC")

	componentPaths = collectComponentPaths(sourceDir, componentPaths)
	fmt.Println("Count", len(componentPaths))

	for _, f := range componentPaths {
		copyDir(f, targetDir+"/"+path.Base(f))
	}

	fmt.Println("Done.")
}

func CopyAllStaticResources() {
	sourceDir := watchPackagePath() + "/staticresources"
	targetDir := defaultPackagePath() + "/staticresources"

	fmt.Println("Start copy static resources")

	copyDir(sourceDir, targetDir)

	fmt.Println("Done.")
}

/**
 * Return
 * 1. Component name
 * 2. Component path
 */
func getComponentNameWithPath(pathStr string) (string, string) {

	if pathStr == watchPackagePath() || pathStr == "/" {
		log.Fatal("cannot detect component name")
	}

	dir := path.Dir(pathStr)
	_, files, _ := getListFiles(dir)

	var name = ""

	for _, f := range files {
		file := path.Base(f)
		if strings.Index(file, ".js-meta.xml") != -1 {
			name = path.Base(path.Dir(f))
		}
	}

	if name != "" {
		return name, dir
	} else {
		return getComponentNameWithPath(dir)
	}
}

func collectComponentPaths(dir string, componentPaths []string) []string {
	lwcDir := dir + "/lwc"

	if !pathIsExist(lwcDir) {
		return append(componentPaths, dir)
	} else {
		_, _, folders := getListFiles(lwcDir)
		for _, folder := range folders {
			componentPaths = collectComponentPaths(folder, componentPaths)
		}
		return componentPaths
	}
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

/**
 * Return
 * 1. List of files and folder
 * 2. List of files
 * 3. List of folder
 */
func getListFiles(path string) ([]string, []string, []string) {
	var paths []string
	var files []string
	var folders []string

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
		paths = append(paths, path+"/"+v.Name())
		if v.IsDir() {
			folders = append(folders, path+"/"+v.Name())
		} else {
			files = append(files, path+"/"+v.Name())
		}
	}

	return paths, files, folders
}

func copyDir(source string, target string) {
	paths, _, _ := getListFiles(source)
	for _, f := range paths {
		stat, _ := os.Stat(f)
		if stat.IsDir() {
			copyDir(f, target+"/"+path.Base(f))
		} else {
			copyFile(f, target+"/"+stat.Name())
		}
	}
}

func copyFile(source string, target string) {
	targetDir := path.Dir(target)
	if !pathIsExist(targetDir) {
		err := os.MkdirAll(targetDir, 0755)
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

func isEqual(sourceFile, targetFile string) bool {
	sf, err := os.Open(sourceFile)
	if err != nil {
		log.Fatal(err)
	}

	if !pathIsExist(targetFile) {
		return false
	}

	df, err := os.Open(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}

	return true
}
