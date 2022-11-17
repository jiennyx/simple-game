package repository

import (
	"context"
	"errors"
)

var (
	ErrInsert = errors.New("insert error")
)

type UserRepository interface {
	Create(ctx context.Context, username, password string) error
}
