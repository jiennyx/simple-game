package logic

import (
	"context"
	"time"

	"google.golang.org/grpc/status"
	"simplegame.com/simplegame/common/api/user"
	"simplegame.com/simplegame/common/clients"
	"simplegame.com/simplegame/common/jwtx"
	"simplegame.com/simplegame/web/model"
	"simplegame.com/simplegame/web/server/errorx"
)

func GetAuth(
	ctx context.Context,
	req model.GetAuthReq,
) (model.GetAuthRsp, error) {
	res := model.GetAuthRsp{}

	rsp, err := clients.UserClient().ExistUser(ctx, &user.ExistUserReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		status, _ := status.FromError(err)
		return res, errorx.ErrorFromStatus(status.Code(), status.Message())
	}
	if !rsp.GetIsExisted() {
		return res, nil
	}

	jwtInfo := jwtx.JwtInfo{
		Username: req.Username,
	}
	res.AuthToken, err = jwtx.GenerateToken(
		jwtx.JWTTypeAuth, jwtInfo, time.Minute*10,
	)
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = jwtx.GenerateToken(
		jwtx.JWTTypeRefresh, jwtInfo, time.Minute*30,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}
