package service

import (
	"errors"
	repo "hita/repository"
	"hita/utils/api"
	"hita/utils/verify"
	"strconv"
)

type ReqSignUp struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Nickname string `form:"nickname" json:"nickname"`
	Gender   string `form:"gender" json:"gender" `
}

type ReqLogIn struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
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

func (req *ReqLogIn) LogIn() (token string, code int, err error) {
	var user = repo.User{
		UserName: req.Username,
	}
	if user.FindUser() == nil {
		if user.Password == req.Password {
			token, err = verify.SignToken(strconv.FormatInt(user.Id, 10))
		} else {

			err = errors.New("wrong password")
			code = api.CodeWrongPassword
		}
	} else {
		err = errors.New("user does not exist")
		code = api.CodeUserNotExist
	}
	return token, code, err
}
