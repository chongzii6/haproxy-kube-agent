package main

import (
	"flag"
	"strings"
)

//CmdCfg keeps global configuration
var CmdCfg = Config{}

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

	if len(*configfile) > 0 {
		CmdCfg.getConf(*configfile)
	}

	if len(*cafile) > 0 {
		CmdCfg.Cafile = *cafile
	}

	if len(*certfile) > 0 {
		CmdCfg.Certfile = *certfile
	}

	if len(*keyfile) > 0 {
		CmdCfg.Keyfile = *keyfile
	}

	if len(*endpoints) > 0 {
		CmdCfg.Endpoints = strings.Split(*endpoints, ",")
	}

	if len(*agentkey) > 0 {
		CmdCfg.Agentkey = *agentkey
	}

	if len(*reqkey) > 0 {
		CmdCfg.Reqkey = *reqkey
	}
}

func main() {
	initConfig()
}
