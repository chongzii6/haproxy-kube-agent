package main

import (
	"github.com/chongzii6/haproxy-kube-agent/agent"
)

func main() {
	agent.InitConfigFromArgs()
	quit := make(chan int)
	agent.EtcdWatch(agent.CmdCfg.Reqkey, quit)
}
