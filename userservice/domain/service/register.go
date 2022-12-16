package service

import (
	"context"

	"simplegame.com/simplegame/common/logx"
	"simplegame.com/simplegame/userservice/domain/entity"
	"simplegame.com/simplegame/userservice/domain/repository"
)

type RegisterDomainService struct {
	userRepo repository.UserRepository

	logger logx.Logger
}

type RegisterConfiguration func(rs *RegisterDomainService)

func NewRegisterDomainService(logger logx.Logger, cfgs ...RegisterConfiguration) (
	RegisterDomainService, error) {
	res := RegisterDomainService{
		logger: logger,
	}
	for _, cfg := range cfgs {
		cfg(&res)
	}

	return res, nil
}

func WithUserRepository(ur repository.UserRepository) RegisterConfiguration {
	return func(rs *RegisterDomainService) {
		rs.userRepo = ur
	}
}

func (service *RegisterDomainService) Register(
	ctx context.Context,
	username, password string,
) error {
	account := entity.Account{
		Username: username,
		Password: password,
	}
	if err := account.Validate(); err != nil {
		return err
	}
	if err := service.userRepo.Create(
		ctx, account.Username, account.Password); err != nil {
		return err
	}
	// TODO
	// register event

	return nil
}
