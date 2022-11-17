package aggregate

import "simplegame.com/simplegame/userservice/domain/entity"

type User struct {
	base    *entity.UserBase
	account *entity.Account
}

func (u *User) GetUid() uint {
	return u.base.Uid
}

func (u *User) GetName() string {
	return u.base.Name
}

func (u *User) GetSex() uint {
	return u.base.Sex
}

func (u *User) GetUsername() string {
	return u.account.Username
}

func (u *User) GetPassword() string {
	return u.account.Password
}

type UserBuilder struct {
	user *User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{user: &User{}}
}

func (builder *UserBuilder) WithUid(uid uint) *UserBuilder {
	builder.user.base.Uid = uid

	return builder
}

func (builder *UserBuilder) WithName(name string) *UserBuilder {
	builder.user.base.Name = name

	return builder
}

func (builder *UserBuilder) WithSex(sex uint) *UserBuilder {
	builder.user.base.Sex = sex

	return builder
}

func (builder *UserBuilder) WithUsername(username string) *UserBuilder {
	builder.user.account.Username = username

	return builder
}

func (builder *UserBuilder) WithPassword(password string) *UserBuilder {
	builder.user.account.Password = password

	return builder
}

func (builder *UserBuilder) Build() *User {
	return builder.user
}
