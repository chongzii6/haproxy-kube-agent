package agent

import (
	"fmt"
	"testing"
)

func TestHandleReq(t *testing.T) {
	key := []byte("/chongzii6/lbreq/testkey")
	value := []byte(`
	{
		"action":"add", 
		"loadbalance_name":"lb1", 
		"target_port":3800, 
		"endpoints":[
			{
				"name":"server1", 
				"ip":"192.168.100.210", 
				"port":38000
			}, 
			{
				"name":"server2", 
				"ip":"192.168.100.211",
				"port":38000
			}
		]
	}`)

	CmdCfg.Cfgpath = "/c/Users/junlinch/vol"
	err := HandleReq(key, value)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetLocalIP(t *testing.T) {
	ip, err := getLocalIP("")
	if err == nil {
		fmt.Println(ip)
	} else {
		fmt.Println(err)
	}

	ip, err = getLocalIP("本地连接")
	if err == nil {
		fmt.Println(ip)
	} else {
		fmt.Println(err)
	}
}
