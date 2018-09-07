package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chongzii6/haproxy-kube-agent/agent"
)

func main() {
	agent.InitConfigFromArgs()
	quit := make(chan int)

	go func() {
		<-time.After(time.Second * 2)

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter Command: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			fmt.Println("command: ", text)

			if text == "quit" {
				quit <- 0
				return
			}
		}
	}()

	agent.EtcdWatch(agent.CmdCfg.GetReqkey(), quit)
}
