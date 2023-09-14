package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type StageDescription struct {
	Number      int    `yaml:"number"`
	Name        string `yaml:"name"`
	Params      string `yaml:"params"`
	Tool        string `yaml:"tool"`
	Description string `yaml:"description"`
	ToolPath    string `yaml:"toolPath"`
}

type Stage struct {
	Stage *StageDescription `yaml:"stage"`
}

type Protocol struct {
	Version   string              `yaml:"version"`
	Author    string              `yaml:"author"`
	Stages    []*StageDescription `yaml:"stages"`
	CreatedAt time.Time           `yaml:"createdAt"`
	UpdatedAt time.Time           `yaml:"updatedAt"`
}

func main() {
	prt := NewProtocol()

	filename, err := filepath.Abs("internal/protocols/test_protocol.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	protocol, err := ParseProtocol(yamlFile, prt)
	log.Println(protocol, err)
	for _, val := range protocol.Stages {
		log.Println(val)
	}
	// log.p
}

func NewProtocol() *Protocol {
	return new(Protocol)
}

func ParseProtocol(path []byte, protocol *Protocol) (*Protocol, error) {
	if err := yaml.Unmarshal(path, protocol); err != nil {
		return nil, err
	}
	return protocol, nil
}
