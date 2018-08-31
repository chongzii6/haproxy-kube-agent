package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
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

var cfg = Config{}

func initConfig() {
	// quiets down kube client library logging
	configfile := flag.String("config", "", "<optional> path to config yaml file")
	cafile := flag.String("ca", "", "ca.pem used for etcd")
	certfile := flag.String("cert", "", "cert.pem used for etcd")
	keyfile := flag.String("key", "", "key.pem used for etcd")
	endpoints := flag.String("endpoints", "", "endpoints of etcd")
	agentkey := flag.String("agent", "", "agent key used in etcd")
	reqkey := flag.String("req", "", "req key used in etcd")
	flag.Parse()

	// f, err := os.Open(*configfile)
	// if err == nil {
	// 	defer f.Close()
	// 	var c Config
	// 	err = gcfg.ReadInto(&c, f)
	// 	if err == nil {
	// 		cfg = c
	// 	} else {
	// 		log.Fatalf("Failed to parse gcfg data: %s", err)
	// 	}
	// }

	if len(*configfile) > 0 {
		cfg.getConf(*configfile)
	}

	if len(*cafile) > 0 {
		cfg.Cafile = *cafile
	}

	if len(*certfile) > 0 {
		cfg.Certfile = *certfile
	}

	if len(*keyfile) > 0 {
		cfg.Keyfile = *keyfile
	}

	if len(*endpoints) > 0 {
		cfg.Endpoints = strings.Split(*endpoints, ",")
	}

	if len(*agentkey) > 0 {
		cfg.Agentkey = *agentkey
	}

	if len(*reqkey) > 0 {
		cfg.Reqkey = *reqkey
	}
}

func main() {
	initConfig()
}
