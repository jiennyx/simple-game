package clients

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type EtcdConfig struct {
	Endpoints         []string
	DialTimeout       int
	DialKeepAliveTime int
}

const (
	keyPrefix = "microservice"
)

func getServicePrefix(service string) string {
	return fmt.Sprintf("%s/%s", keyPrefix, service)
}

func getAddr(key string) string {
	sli := strings.Split(key, "/")
	if len(sli) < 3 {
		return ""
	}

	return sli[2]
}

func RegisterService(name, ip string, port int, conf EtcdConfig) error {
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
	key := fmt.Sprintf("%s/%s:%d", getServicePrefix(name), ip, port)
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

func DiscoverService(services []string, conf EtcdConfig) (*Pool, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:         conf.Endpoints,
		DialTimeout:       time.Duration(conf.DialTimeout) * time.Second,
		DialKeepAliveTime: time.Duration(conf.DialKeepAliveTime) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	res := &Pool{
		clients: make(map[string]map[string][]any),
		lock:    sync.RWMutex{},
	}

	ctx := context.Background()
	for _, service := range services {
		go handleServicePool(ctx, client, res, service)
	}

	return res, nil
}

func handleServicePool(
	ctx context.Context,
	client *clientv3.Client,
	pool *Pool,
	service string,
) {
	pre := getServicePrefix(service)
	getRes, err := client.Get(ctx, pre, clientv3.WithPrefix())
	if err != nil {
		zap.L().Error(
			"get service with prefix failed",
			zap.String("prefix", pre),
		)
	}
	for _, kv := range getRes.Kvs {
		add(pool, kv)
	}
	ch := client.Watch(ctx, pre,
		clientv3.WithPrefix(), clientv3.WithPrevKV())
	for item := range ch {
		for _, event := range item.Events {
			switch event.Type {
			case mvccpb.PUT:
				add(pool, event.Kv)
			case mvccpb.DELETE:
				del(pool, event.Kv)
			}
		}
	}
}

// TODO
func add(pool *Pool, kv *mvccpb.KeyValue) {
	key, value := string(kv.Key), string(kv.Value)
	sli := strings.Split(key, "/")
	if len(sli) < 3 {
		return
	}
	service, addr := sli[1], sli[2]
	pool.lock.Lock()
	if _, ok := pool.clients[service]; !ok {
		pool.clients[service] = make(map[string][]any)
	}
	if _, ok := pool.clients[service][value]; !ok {
		pool.clients[service][value] = make([]any, 0)
	}
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		zap.L().Error(
			"get service with prefix failed",
			zap.String("service", service),
		)
	}
	pool.clients[service][value] =
		append(pool.clients[service][value], newClient(service, conn))
	pool.lock.Unlock()
}

func del(pool *Pool, kv *mvccpb.KeyValue) {
	key, value := string(kv.Key), string(kv.Value)
	sli := strings.Split(key, "/")
	if len(sli) < 3 {
		return
	}
	service, addr := sli[1], sli[2]
	pool.lock.Lock()
	if _, ok := pool.clients[service]; !ok {
		return
	}
	if _, ok := pool.clients[service][value]; !ok {
		return
	}
	for i, v := range pool.clients[service][value] {
		if v == addr {
			pool.clients[service][value] = append(
				pool.clients[service][value][:i],
				pool.clients[service][value][i+1:]...,
			)
		}
	}
	pool.lock.Unlock()
}

func CancelService(service, ip string, port int, conf EtcdConfig) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:         conf.Endpoints,
		DialTimeout:       time.Duration(conf.DialTimeout) * time.Second,
		DialKeepAliveTime: time.Duration(conf.DialKeepAliveTime) * time.Second,
	})
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s:%d", getServicePrefix(service), ip, port)
	ctx := context.Background()
	_, err = client.Delete(ctx, key)

	return err
}
