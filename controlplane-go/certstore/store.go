package certstore

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func SaveCerts(cli *clientv3.Client, cpName string, ca, crt, key []byte) error {
	pairs := map[string][]byte{
		fmt.Sprintf("/controlplane/%s/certs/ca.crt", cpName):     ca,
		fmt.Sprintf("/controlplane/%s/certs/server.crt", cpName): crt,
		fmt.Sprintf("/controlplane/%s/certs/server.key", cpName): key,
	}

	for k, v := range pairs {
		_, err := cli.Put(context.Background(), k, string(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadCerts(cli *clientv3.Client, cpName string) (ca, crt, key []byte, err error) {
	read := func(name string) ([]byte, error) {
		resp, err := cli.Get(context.Background(), fmt.Sprintf("/controlplane/%s/certs/%s", cpName, name))
		if err != nil || len(resp.Kvs) == 0 {
			return nil, fmt.Errorf("missing cert %s", name)
		}
		return resp.Kvs[0].Value, nil
	}
	ca, err = read("ca.crt")
	if err != nil {
		return
	}
	crt, err = read("server.crt")
	if err != nil {
		return
	}
	key, err = read("server.key")
	if err != nil {
		return
	}
	return
}
