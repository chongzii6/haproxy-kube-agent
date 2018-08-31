package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

//Config global config
type Config struct {
	Cafile    string
	Certfile  string
	Keyfile   string
	Endpoints []string
	Agentkey  string
	Reqkey    string
}

func (c *Config) getConf(f string) {

	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
