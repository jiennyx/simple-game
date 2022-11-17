package service

import (
	"context"

	"simplegame.com/simplegame/userservice/domain/aggregate"
	"simplegame.com/simplegame/userservice/domain/repository"
	"simplegame.com/simplegame/userservice/domain/service"
)

type RegisterApplicationService struct {
	userRepo              repository.UserRepository
	registerDomainService service.RegisterDomainService
}

type RegisterConfiguration func(rs *RegisterApplicationService) error

func NewRegisterApplicationService(cfgs ...RegisterConfiguration) (
	*RegisterApplicationService, error) {
	res := &RegisterApplicationService{}
	for _, cfg := range cfgs {
		err := cfg(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func WithUserRepository(ur repository.UserRepository) RegisterConfiguration {
	return func(rs *RegisterApplicationService) error {
		rs.userRepo = ur

		return nil
	}
}

func WithRegisterDomainService(
	domainService service.RegisterDomainService,
) RegisterConfiguration {
	return func(rs *RegisterApplicationService) error {
		rs.registerDomainService = domainService

		return nil
	}
}

func (service *RegisterApplicationService) Register(
	ctx context.Context,
	username, password string,
) error {
	return service.registerDomainService.Register(ctx, username, password)
}

func (service *RegisterApplicationService) GetUser(
	ctx context.Context,
	uid uint,
) (*aggregate.User, error) {
	return service.userRepo.GetByUid(ctx, uid)
}
