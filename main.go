package main

import (
	"flag"
	"os"
	"strings"

	"gopkg.in/gcfg.v1"
)

//Config global config
type Config struct {
	cafile    string
	certfile  string
	keyfile   string
	endpoints []string
	agentkey  string
	reqkey    string
}

var cfg = Config{}

func initConfig() {
	// quiets down kube client library logging
	configfile := flag.String("config", "ageng.cfg", "<optional> path to config file")
	cafile := flag.String("ca", "", "ca.pem used for etcd")
	certfile := flag.String("cert", "", "cert.pem used for etcd")
	keyfile := flag.String("key", "", "key.pem used for etcd")
	endpoints := flag.String("endpoints", "", "endpoints of etcd")
	agentkey := flag.String("agent", "", "agent key used in etcd")
	reqkey := flag.String("req", "", "req key used in etcd")
	flag.Parse()

	f, err := os.Open(*configfile)
	if err == nil {
		defer f.Close()
		var c Config
		err = gcfg.ReadInto(&c, f)
		if err == nil {
			cfg = c
		}
	}

	if len(*cafile) > 0 {
		cfg.cafile = *cafile
	}

	if len(*certfile) > 0 {
		cfg.certfile = *certfile
	}

	if len(*keyfile) > 0 {
		cfg.keyfile = *keyfile
	}

	if len(*endpoints) > 0 {
		cfg.endpoints = strings.Split(*endpoints, ",")
	}

	if len(*agentkey) > 0 {
		cfg.agentkey = *agentkey
	}

	if len(*reqkey) > 0 {
		cfg.reqkey = *reqkey
	}
}

func main() {
	initConfig()
}
