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

	etcdWatch(CmdCfg.Reqkey, quit)
}
