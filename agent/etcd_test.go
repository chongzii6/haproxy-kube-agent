package agent

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEtcdWatch(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	quit := make(chan int)
	go func() {
		<-time.After(time.Minute * 2)

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Enter text: ")
		// text, _ := reader.ReadString('\n')
		// fmt.Println(text)

		quit <- 0
	}()

	EtcdWatch(CmdCfg.GetReqkey(), quit)
}

func TestWatchAndPut(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	quit := make(chan int)
	go func() {
		<-time.After(time.Second * 2)

		req := &Request{
			Action:     ADD,
			LbName:     "lb2",
			TargetPort: 3900,
			Endpoints: []Endpoint{
				{Name: "server1", IP: "192.168.100.210", Port: 38000},
				{Name: "server2", IP: "192.168.100.211", Port: 38000},
			},
		}
		SendReq(req)

		<-time.After(time.Second * 15)
		quit <- 0
	}()

	EtcdWatch(CmdCfg.GetReqkey(), quit)
}

func TestGet(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	key := "/chongzii6/b"
	str, err := EtcdGet(key)
	if err == nil && str != "" {
		fmt.Println(str)
	}

	key = "/chongzii6/a"
	str, err = EtcdGet(key)
	if err == nil && str != "" {
		fmt.Println(str)
	}
}

func TestWatchAndDel(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	quit := make(chan int)
	go func() {
		<-time.After(time.Second * 2)

		req := &Request{
			Action: DELETE,
			LbName: "lb2",
		}
		SendReq(req)

		<-time.After(time.Second * 15)
		quit <- 0
	}()

	EtcdWatch(CmdCfg.GetReqkey(), quit)
}
