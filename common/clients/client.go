package clients

import (
	"log"

	"google.golang.org/grpc"
	"simplegame.com/simplegame/common/api/user"
)

var (
	userClient user.UserClient
)

// TODO
func UserClient() user.UserClient {
	if userClient == nil {
		conn, err := grpc.Dial(
			"dns:///simple-game-userservice:80",
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(
				`{"loadBalancingPolicy": "round_robin"}`,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		userClient = user.NewUserClient(conn)
	}

	return userClient
}
