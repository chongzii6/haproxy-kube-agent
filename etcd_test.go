package main

import (
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	etcdWatch(CmdCfg.Reqkey, time.After(time.Second*20))
}
