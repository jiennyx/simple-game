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
		conn, err := grpc.Dial("127.0.0.1:30001", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		userClient = user.NewUserClient(conn)
	}

	return userClient
}
