package agent

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/satori/go.uuid"
)

//RequestType define request
type RequestType int32

//ADD loadbalancer
const (
	ADD    RequestType = 1
	UPDATE RequestType = 2
	DELETE RequestType = 3
)

//Request struct
type Request struct {
	Action     RequestType `json:"action"`
	LbName     string      `json:"loadbalance_name"`
	TargetPort int         `json:"target_port,omitempty"`
	Endpoints  []Endpoint  `json:"endpoints,omitempty"`
	SvcName    string      `json:"service_name,omitempty"`
}

//LBState store loadbalancer status
type LBState struct {
	HostIP      string `json:"ip"`
	Port        string `json:"port,omitempty"`
	ContainerID string `json:"cid"`
}

var hostIP string

//GetHostIP return host ip
func GetHostIP() string {
	if hostIP == "" {
		hostIP, _ = getLocalIP(CmdCfg.Ifname)
	}

	return hostIP
}

//SendReq put request to etcd
func SendReq(req *Request) error {
	u1 := uuid.NewV4()
	key := fmt.Sprintf("%s/%s", CmdCfg.GetReqkey(), u1)
	by, err := json.Marshal(req)
	if err == nil {
		EtcdPut(key, string(by))
	}

	return err
}

//HandleReq handle watched request
func HandleReq(reqKey []byte, reqVal []byte) error {
	var req Request
	key := string(reqKey)
	isAny := CmdCfg.IsAnyKey(key)

	if isAny {

	}

	err := json.Unmarshal(reqVal, &req)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	switch req.Action {
	case ADD:
		err = addLoadBalancer(req.LbName, req.Endpoints, req.TargetPort)
	case UPDATE:
		err = updateLoadBalancer(req.LbName, req.Endpoints)
	case DELETE:
		err = deleteLoadBalancer(req.LbName)
	}
	if err == nil {
		EtcdDel(key)
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
	cid, err := RunHaproxy(name, dockerpath, p, true)
	if err != nil {
		return err
	}

	lbkey := fmt.Sprintf("%s/%s", CmdCfg.Agentkey, name)

	st := &LBState{
		HostIP:      CmdCfg.PublicIP,
		ContainerID: cid,
		Port:        p}

	text, err := json.Marshal(st)
	if err == nil {
		EtcdPut(lbkey, string(text))
	}

	return err
}

func updateLoadBalancer(name string, endpoints []Endpoint) error {
	return nil
}

func deleteLoadBalancer(name string) error {
	lbkey := fmt.Sprintf("%s/%s", CmdCfg.Agentkey, name)
	text, err := EtcdGet(lbkey)
	if err != nil {
		return err
	}

	st := &LBState{}
	err = json.Unmarshal([]byte(text), st)
	if err == nil {
		err = DelHaproxy(st.ContainerID)
		if err != nil {
			return err
		}
	}

	err = EtcdDel(lbkey)
	return err
}

func getLocalIP(ifname string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		if ifname != "" && ifname != i.Name {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.To4()
			case *net.IPAddr:
				ip = v.IP.To4()
			}

			if ip != nil {
				// fmt.Println(i.Name, ip)
				if !ip.IsLoopback() {
					return ip.String(), nil
				}
			}
		}
	}

	return "", nil
}
