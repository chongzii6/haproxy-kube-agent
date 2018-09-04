package agent

import (
	"testing"
)

func TestGetHaList(t *testing.T) {
	RunHaproxy("haproxy1", "/c/Users/junlinch/vol/haproxy.cfg", "2800", true)
	GetHaList()
}
