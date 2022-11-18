package facade

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"simplegame.com/simplegame/common/api/user"
	"simplegame.com/simplegame/userservice/application/service"
)

type UserServer struct {
	registerService service.RegisterApplicationService
	user.UnimplementedUserServer
}

type UserConfiguration func(us *UserServer) error

func NewUserServer(cfgs ...UserConfiguration) (UserServer, error) {
	res := UserServer{}
	for _, cfg := range cfgs {
		err := cfg(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func WithRegisterApplicationService(
	applicationService service.RegisterApplicationService,
) UserConfiguration {
	return func(us *UserServer) error {
		us.registerService = applicationService

		return nil
	}
}

func (server *UserServer) Register(
	ctx context.Context,
	req *user.RegisterReq,
) (*user.RegisterRsp, error) {
	return nil, nil
}

func (server *UserServer) ExistUser(
	ctx context.Context,
	req *user.ExistUserReq,
) (*user.ExistUserRsp, error) {
	res := &user.ExistUserRsp{}
	isExisted, err := server.registerService.
		ExistUser(ctx, req.Username, req.Password)
	if err != nil {
		return res, status.Error(codes.Internal, codes.Internal.String())
	}

	res.IsExisted = isExisted

	return res, nil
}
