package service

import (
	repo "hita/repository"
)

type ReqSignUp struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Nickname string `form:"nickname" json:"nickname"`
	Gender string `form:"gender" json:"gender" `
}

func (req* ReqSignUp) SignUp() (id int64, err error) {
	user:= repo.User{
		UserName: req.Username,
		Nickname: req.Nickname,
		Gender: req.Gender,
		Password: req.Password,
	}
	return user.AddUser()
}