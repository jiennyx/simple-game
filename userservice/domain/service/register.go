package service

import (
	"context"

	"simplegame.com/simplegame/userservice/domain/entity"
	"simplegame.com/simplegame/userservice/domain/repository"
)

type RegisterDomainService struct {
	userRepo repository.UserRepository
}

type RegisterConfiguration func(rs *RegisterDomainService) error

func NewRegisterDomainService(cfgs ...RegisterConfiguration) (
	RegisterDomainService, error) {
	res := RegisterDomainService{}
	for _, cfg := range cfgs {
		err := cfg(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func WithUserRepository(ur repository.UserRepository) RegisterConfiguration {
	return func(rs *RegisterDomainService) error {
		rs.userRepo = ur

		return nil
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
