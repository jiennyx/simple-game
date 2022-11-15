package repository

import (
	"context"

	"simplegame.com/simplegame/userservice/domain/entity"
)

type AccountRepository interface {
	Save(ctx context.Context, account *entity.Account) error
}
