package facade

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"simplegame.com/simplegame/common/api/user"
	"simplegame.com/simplegame/common/logx"
	"simplegame.com/simplegame/userservice/application/service"
)

type UserServer struct {
	registerService service.RegisterApplicationService

	logger logx.Logger

	user.UnimplementedUserServer
}

type UserConfiguration func(us *UserServer)

func NewUserServer(logger logx.Logger, cfgs ...UserConfiguration) (UserServer, error) {
	res := UserServer{
		logger: logger,
	}
	for _, cfg := range cfgs {
		cfg(&res)
	}

	return res, nil
}

func WithRegisterApplicationService(
	applicationService service.RegisterApplicationService,
) UserConfiguration {
	return func(us *UserServer) {
		us.registerService = applicationService
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
