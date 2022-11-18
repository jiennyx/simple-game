package repository

import (
	"context"
	"errors"

	"simplegame.com/simplegame/userservice/domain/aggregate"
)

var (
	ErrInsert       = errors.New("insert error")
	ErrUserNotFound = errors.New("user not found")
	ErrInternal     = errors.New("internal error")
)

type UserRepository interface {
	Create(ctx context.Context, username, password string) error
	GetByUid(ctx context.Context, uid uint) (*aggregate.User, error)
	ExistUser(ctx context.Context, username, password string) (bool, error)
}
