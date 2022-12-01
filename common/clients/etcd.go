package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.etcd.io/etcd/client/v3"
	"simplegame.com/simplegame/common/netx"
)

type EtcdConfig struct {
	Endpoints         []string
	DialTimeout       int
	DialKeepAliveTime int
}

const (
	keyPrefix = "/microservice"
)

func RegisterService(name string, port int, conf EtcdConfig) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:         conf.Endpoints,
		DialTimeout:       time.Duration(conf.DialTimeout) * time.Second,
		DialKeepAliveTime: time.Duration(conf.DialKeepAliveTime) * time.Second,
	})
	if err != nil {
		return err
	}
	color := os.Getenv("color")
	if color == "" {
		color = "default"
	}
	ip, err := netx.GetLocalIP()
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s/%s:%d", keyPrefix, name, ip, port)
	ctx := context.Background()
	lease, err := client.Grant(ctx, 3)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, key, color, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	_, err = client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	return nil
}
