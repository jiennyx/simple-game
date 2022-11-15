package dao

import (
	"context"

	"gorm.io/gorm"
	"simplegame.com/simplegame/userservice/domain/entity"
	"simplegame.com/simplegame/userservice/domain/repository"
)

type accountRepo struct {
	db *gorm.DB
}

var _ repository.AccountRepository = (*accountRepo)(nil)

func NewAccountRepository(db *gorm.DB) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

func (ar *accountRepo) Save(ctx context.Context, account *entity.Account) error {
	return nil
}
