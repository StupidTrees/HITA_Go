package service

import (
	"errors"
	repo "hita/repository"
	"hita/utils/api"
	"hita/utils/security"
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

func (req *ReqSignUp) SignUp() (id int64, token string, publicKey string, code int, err error) {
	var user = repo.User{
		UserName: req.Username,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Password: req.Password,
	}
	//生成用户公私钥对
	user.PublicKey, user.PrivateKey, err = security.GenerateRSAKeysStr()
	if err != nil {
		publicKey = ""
		code = api.CodeOtherError
		return
	}
	//密码使用用户私钥加密存储
	user.Password = security.EncryptWithPrivateKey(user.Password, user.PrivateKey)
	id, err = user.AddUser()
	if err == nil {
		token, err = verify.SignToken(strconv.FormatInt(id, 10))
		publicKey = user.PublicKey
	} else {
		code = api.CodeUserExists
		err = errors.New("user already exists")
	}
	return
}

func (req *ReqLogIn) LogIn() (token string, code int, err error) {
	var user = repo.User{
		UserName: req.Username,
	}
	if user.FindUser() == nil {
		realPassword := security.DecryptWithPublicKey(user.Password, user.PublicKey)
		if realPassword == req.Password { //} security.DecryptWithPrivateKey(req.Password,user.PrivateKey) {
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
