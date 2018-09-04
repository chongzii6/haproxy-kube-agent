package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//Request struct
type Request struct {
	Action     string     `json:"action"`
	LbName     string     `json:"loadbalance_name"`
	TargetPort int        `json:"target_port,omitempty"`
	Endpoints  []Endpoint `json:"endpoints,omitempty"`
}

//HandleReq handle watched request
func HandleReq(reqKey []byte, reqVal []byte) error {
	var req Request
	err := json.Unmarshal(reqVal, &req)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	switch req.Action {
	case "add":
		err = addLoadBalancer(req.LbName, req.Endpoints, req.TargetPort)
	case "upd":
		updateLoadBalancer(req.LbName, req.Endpoints)
	case "del":
		deleteLoadBalancer(req.LbName)
	}

	return err
}

func addLoadBalancer(name string, endpoints []Endpoint, port int) error {
	dockerpath := fmt.Sprintf("%s/%s-haproxy.cfg", CmdCfg.Cfgpath, name)
	cfgpath := dockerpath

	goos := runtime.GOOS
	if goos == "windows" {
		plist := strings.Split(dockerpath, "/")
		cfgpath = fmt.Sprintf("%s:%c%s", plist[1], filepath.Separator, strings.Join(plist[2:], string(filepath.Separator)))
	}

	file, err := os.Create(cfgpath)
	if err != nil {
		return err
	}

	WriteHaproxyCfg(file, endpoints)
	file.Close()

	if port == 0 {

	}
	p := fmt.Sprintf("%d", port)
	err = RunHaproxy(name, dockerpath, p, true)
	if err != nil {
		return err
	}
	return nil
}

func updateLoadBalancer(Name string, Endpoints []Endpoint) {

}

func deleteLoadBalancer(Name string) {

}
