package agent

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

const (
	dialTimeout = 5 * time.Second
)

func newClientCfg() (*clientv3.Config, error) {
	var cfgtls *transport.TLSInfo
	tlsinfo := transport.TLSInfo{}

	cfg := &clientv3.Config{
		Endpoints:   CmdCfg.Endpoints,
		DialTimeout: dialTimeout,
	}

	if CmdCfg.Certfile != "" {
		tlsinfo.CertFile = CmdCfg.Certfile
		cfgtls = &tlsinfo
	}

	if CmdCfg.Keyfile != "" {
		tlsinfo.KeyFile = CmdCfg.Keyfile
		cfgtls = &tlsinfo
	}

	if CmdCfg.Cafile != "" {
		tlsinfo.TrustedCAFile = CmdCfg.Cafile
		cfgtls = &tlsinfo
	}

	if cfgtls != nil {
		clientTLS, err := cfgtls.ClientConfig()
		if err != nil {
			return nil, err
		}
		cfg.TLS = clientTLS
	}
	return cfg, nil
}

func etcdWatch(key string, quit chan int) error {
	cfg, err := newClientCfg()
	if err != nil {
		return err
	}

	client, err := clientv3.New(*cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	wc := client.Watch(context.Background(), key, clientv3.WithPrefix())

	select {
	case <-quit:
		return nil
	case resp := <-wc:
		for _, e := range resp.Events {
			fmt.Printf("%s", e)
		}
	}
	return nil
}

func etcdGet() error {
	cfg, err := newClientCfg()
	if err != nil {
		return err
	}

	client, err := clientv3.New(*cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	return nil
}
