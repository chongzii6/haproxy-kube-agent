package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chongzii6/haproxy-kube-agent/agent"
)

func main() {
	agent.InitConfigFromArgs()
	quit := make(chan int)

	c := make(chan os.Signal)
	signal.Notify(c)

	go func() {
		// <-time.After(time.Second * 2)
		s := <-c
		log.Println("Got signal:", s) //Got signal: terminated
		quit <- 0

		// for {
		// 	reader := bufio.NewReader(os.Stdin)
		// 	fmt.Print("Enter Command: ")
		// 	text, err := reader.ReadString('\n')
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		break
		// 	}
		// 	text = strings.TrimSpace(text)
		// 	fmt.Println("command: ", text)

		// 	if text == "quit" {
		// 		quit <- 0
		// 		return
		// 	}
		// }
	}()

	// go func() {
	// 	agent.EtcdWatch(agent.CmdCfg.GetAnyReqkey(), quit)
	// }()

	agent.EtcdWatch(agent.CmdCfg.GetReqkey(), quit)
}
