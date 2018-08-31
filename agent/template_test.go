package agent

import (
	"os"
	"testing"
)

func TestWriteHaproxyCfg(t *testing.T) {
	tests := []Endpoint{
		{Name: "server1", IP: "192.168.100.200", Port: 80},
		{Name: "server2", IP: "192.168.100.201", Port: 80},
		{Name: "server3", IP: "192.168.100.202", Port: 80},
	}

	WriteHaproxyCfg(os.Stdout, tests)
}
