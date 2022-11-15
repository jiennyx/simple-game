package service

import "simplegame.com/simplegame/userservice/domain/repository"

type RegisterService struct {
	accountRepo repository.AccountRepository
}

type RegisterConfiguration func(as *RegisterService) error

func NewRegisterService(cfgs ...RegisterConfiguration) (*RegisterService, error) {
	res := &RegisterService{}
	for _, cfg := range cfgs {
		err := cfg(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func WithAccountRepository(ar repository.AccountRepository) RegisterConfiguration {
	return func(as *RegisterService) error {
		as.accountRepo = ar

		return nil
	}
}
