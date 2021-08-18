package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type SfdxProject struct {
	PackageDirectories []struct {
		Path    string `json:"path"`
		Default bool   `json:"default"`
		Watch   bool   `json:"watch"`
	} `json:"packageDirectories"`

	Namespace        string `json:"namespace"`
	SfdcLoginUrl     string `json:"sfdcLoginUrl"`
	SourceApiVersion string `json:"sourceApiVersion"`
}

func LoadConfig(dir string) SfdxProject {

	config := SfdxProject{}

	jsonFile, err := os.Open(dir + "/sfdx-project.json")
	if err != nil {
		log.Fatal("Error open sfdx-project.json file")
	}

	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error read sfdx-project.json file")
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("Error parse sfdx-project.json file")
	}

	return config
}

func (s *SfdxProject) GetWatchPackage() string {
	for _, p := range s.PackageDirectories {
		if p.Watch {
			return p.Path
		}
	}

	return ""
}

func (s *SfdxProject) GetDefaultPackage() string {
	for _, p := range s.PackageDirectories {
		if p.Default {
			return p.Path
		}
	}

	return ""
}