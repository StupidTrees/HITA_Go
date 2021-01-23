package service

import (
	"hita/lib/verify"
	repo "hita/repository"
	"strconv"
)

type ReqSignUp struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Nickname string `form:"nickname" json:"nickname"`
	Gender   string `form:"gender" json:"gender" `
}

func (req *ReqSignUp) SignUp() (id int64, token string, err error) {
	var user = repo.User{
		UserName: req.Username,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Password: req.Password,
	}
	id, err = user.AddUser()
	if err == nil {
		token, err = verify.SignToken(strconv.FormatInt(id, 10))
	}
	return id, token, err
}
