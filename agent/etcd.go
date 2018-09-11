package agent

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
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

func newClient() (*clientv3.Client, error) {
	cfg, err := newClientCfg()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := clientv3.New(*cfg)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

//EtcdWatch watch key
func EtcdWatch(key string, quit chan int) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	wc := client.Watch(context.Background(), key, clientv3.WithPrefix())

	log.Printf("watching: %s\n", key)
	for {
		select {
		case <-quit:
			log.Println("quit")
			return nil
		case resp := <-wc:
			for _, e := range resp.Events {
				log.Printf("%s key:%s, value:%s\n", e.Type, e.Kv.Key, e.Kv.Value)
				if e.Type == mvccpb.PUT {
					err = HandleReq(e.Kv.Key, e.Kv.Value)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

//EtcdGet get
func EtcdGet(key string) (string, error) {
	client, err := newClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	var resp *clientv3.GetResponse
	if resp, err = client.Get(context.Background(), key); err != nil {
		log.Fatalln(err)
		return "", err
	}
	log.Println("resp: ", resp)

	var ret string
	for _, kv := range resp.Kvs {
		ret += string(kv.Value)
	}
	return ret, nil
}

//EtcdPut put
func EtcdPut(key string, val string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	var resp *clientv3.PutResponse
	if resp, err = client.Put(context.Background(), key, val); err != nil {
		log.Fatalln(err)
		return err
	}

	log.Println("resp: ", resp)
	return nil
}

//EtcdDel del key
func EtcdDel(key string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	var resp *clientv3.DeleteResponse
	if resp, err = client.Delete(context.Background(), key); err != nil {
		log.Fatalln(err)
		return nil
	}

	log.Println("resp: ", resp)
	return nil
}

//EtcdMutex lock
func EtcdMutex() error {
	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	s1, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock/")

}
