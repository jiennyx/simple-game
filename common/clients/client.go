package clients

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"simplegame.com/simplegame/common/api/user"
)

var pool *Pool

type Pool struct {
	clients map[string]map[string][]any
	lock    sync.RWMutex
}

const (
	Userservice = "userservice"
)

func User(ctx context.Context) user.UserClient {
	return getClient(Userservice, getColor(ctx)).(user.UserClient)
}

func getColor(ctx context.Context) string {
	color, ok := ctx.Value("color").(string)
	if !ok || color == "" {
		color = "default"
	}

	return color
}

func getClient(
	service, color string,
) any {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	servicePool, ok := pool.clients[service]
	if !ok {
		return nil
	}
	clients, ok := servicePool[color]
	if !ok {
		return nil
	}

	rand.Seed(time.Now().Unix())

	return clients[rand.Intn(len(clients))]
}

func newClient(service string, conn grpc.ClientConnInterface) interface{} {
	switch service {
	case Userservice:
		return user.NewUserClient(conn)
	}

	return nil
}
