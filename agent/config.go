package agent

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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
	Cfgpath   string
	Ifname    string
	AnyIP     string
	PublicIP  string
}

func (c *Config) getConf(f string) {
	// if !filepath.IsAbs(f) {
	// 	abspath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	f = filepath.Join(abspath, f)
	// }

	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

//GetReqkey return reqkey with ip
func (c *Config) GetReqkey() string {
	return fmt.Sprintf("%s/%s", c.Reqkey, GetHostIP())
}

var anyIP string

//GetAnyReqkey return reqkey for any
func (c *Config) GetAnyReqkey() string {
	if anyIP == "" {
		anyIP = fmt.Sprintf("%s/%s", c.Reqkey, c.AnyIP)
	}
	return anyIP
}

//IsAnyKey check any key
func (c *Config) IsAnyKey(key string) bool {
	return strings.HasPrefix(key, c.GetAnyReqkey())
}

//CmdCfg keeps global configuration
var CmdCfg = Config{}

//InitConfigFromArgs handle command-line options
func InitConfigFromArgs() {
	// quiets down kube client library logging
	configfile := flag.String("config", "", "<optional> path to config yaml file")
	cafile := flag.String("ca", "", "ca.pem used for etcd")
	certfile := flag.String("cert", "", "cert.pem used for etcd")
	keyfile := flag.String("key", "", "key.pem used for etcd")
	endpoints := flag.String("endpoints", "", "endpoints of etcd")
	agentkey := flag.String("agent", "", "agent key used in etcd")
	reqkey := flag.String("req", "", "req key used in etcd")
	cfgpath := flag.String("path", "", "req key used in etcd")
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

	if len(*cfgpath) > 0 {
		CmdCfg.Cfgpath = *cfgpath
	}
}
