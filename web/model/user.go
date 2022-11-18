package model

type RegisterReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterRsp struct{}

type GetAuthReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type GetAuthRsp struct {
	AuthToken    string `json:"authToken"`
	RefreshToken string `json:"refreshToken"`
}
