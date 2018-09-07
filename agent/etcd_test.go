package agent

import (
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

	EtcdWatch(CmdCfg.Reqkey, quit)
}

func TestWatchAndPut(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	quit := make(chan int)
	go func() {
		<-time.After(time.Second * 2)

		// key := CmdCfg.Reqkey + "/testkey"
		// value := `
		// {
		// 	"action":"add",
		// 	"loadbalance_name":"lb1",
		// 	"target_port":3800,
		// 	"endpoints":[
		// 		{
		// 			"name":"server1",
		// 			"ip":"192.168.100.210",
		// 			"port":38000
		// 		},
		// 		{
		// 			"name":"server2",
		// 			"ip":"192.168.100.211",
		// 			"port":38000
		// 		}
		// 	]
		// }`

		req := &Request{
			Action:     "add",
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

	EtcdWatch(CmdCfg.Reqkey, quit)
}

func TestWatchAndDel(t *testing.T) {
	os.Chdir("..")
	CmdCfg.getConf("agent.yml")

	quit := make(chan int)
	go func() {
		<-time.After(time.Second * 2)

		// key := CmdCfg.Reqkey + "/testkey"
		// value := `
		// {
		// 	"action":"del",
		// 	"loadbalance_name":"lb1"
		// }`
		// EtcdPut(key, value)

		req := &Request{
			Action: "del",
			LbName: "lb2",
		}
		SendReq(req)

		<-time.After(time.Second * 15)
		quit <- 0
	}()

	EtcdWatch(CmdCfg.Reqkey, quit)
}
