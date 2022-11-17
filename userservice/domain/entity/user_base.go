package entity

type UserBase struct {
	Uid  uint
	Name string
	Sex  uint // 0-
}

func (u *UserBase) Validate() error {
	// TODO
	return nil
}
