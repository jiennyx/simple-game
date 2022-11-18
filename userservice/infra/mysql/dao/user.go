package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"simplegame.com/simplegame/userservice/domain/aggregate"
	"simplegame.com/simplegame/userservice/domain/repository"
)

type userRepo struct {
	db *gorm.DB
}

var _ repository.UserRepository = (*userRepo)(nil)

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

type userPO struct {
	gorm.Model

	Uid          uint
	Name         string
	Sex          uint
	Username     string
	Password     string
	RegisterTime time.Time
}

func (u *userPO) TableName() string {
	return "user"
}

func (u *userPO) toAggregate() *aggregate.User {
	return aggregate.NewUserBuilder().
		WithUid(u.Uid).
		WithName(u.Name).
		WithSex(u.Sex).
		WithUsername(u.Username).
		WithPassword(u.Password).
		Build()
}

func newFromAggregate(user *aggregate.User) *userPO {
	return &userPO{
		Username: user.GetUsername(),
		Password: user.GetPassword(),
		Uid:      user.GetUid(),
		Name:     user.GetName(),
		Sex:      user.GetSex(),
	}
}

func (repo *userRepo) Create(
	ctx context.Context,
	username, password string,
) error {
	if repo.db.WithContext(ctx).Create(&userPO{
		Username: username,
		Password: password,
	}).Error != nil {
		return repository.ErrInsert
	}

	return nil
}

func (repo *userRepo) GetByUid(
	ctx context.Context,
	uid uint,
) (*aggregate.User, error) {
	po := new(userPO)
	err := repo.db.WithContext(ctx).Where("uid = ?", uid).Take(po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrUserNotFound
		}

		return nil, repository.ErrInternal
	}

	return po.toAggregate(), nil
}

func (repo *userRepo) ExistUser(
	ctx context.Context,
	username, password string,
) (bool, error) {
	po := new(userPO)
	err := repo.db.WithContext(ctx).
		Where("username = ? and password = ?", username, password).
		First(po).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, repository.ErrUserNotFound
		}

		return false, repository.ErrInternal
	}

	return true, nil
}
