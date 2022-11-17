package entity

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	md5Sault = "sault"
)

type Account struct {
	Username string
	Password string // md5
}

func NewAccount(username, password string) *Account {
	return &Account{
		Username: username,
		Password: password,
	}
}

func (u *Account) Validate() error {
	// TODO
	return nil
}

func (u *Account) MD5Password() error {
	m5 := md5.New()
	_, err := m5.Write([]byte(u.Password))
	if err != nil {
		return err
	}
	_, err = m5.Write([]byte(md5Sault))
	if err != nil {
		return err
	}

	u.Password = hex.EncodeToString(m5.Sum(nil))

	return nil
}
