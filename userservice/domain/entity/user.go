package entity

type Account struct {
	Uid      uint64
	Username string
	Password string // md5
}
