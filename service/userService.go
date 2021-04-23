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

type RespLogin struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Avatar    int64  `json:"avatar"`
	Signature string `json:"signature"`
	StudentId string `json:"studentId"`
	School    string `json:"school"`
	PublicKey string `json:"publicKey"`
	Token     string `json:"token"`
}

func (req *ReqSignUp) SignUp() (data RespLogin, code int, err error) {
	var user = repo.User{
		UserName: req.Username,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Password: req.Password,
	}
	//生成用户公私钥对
	user.PublicKey, user.PrivateKey, err = security.GenerateRSAKeysStr()
	if err != nil {
		code = api.CodeOtherError
		return
	}
	//密码使用用户私钥加密存储
	user.Password = security.EncryptWithPrivateKey(user.Password, user.PrivateKey)
	_, err = user.AddUser()
	if err == nil {
		data.Token, err = verify.SignToken(strconv.FormatInt(user.Id, 10))
		data.PublicKey = user.PublicKey
		data.Username = user.UserName
		data.Gender = user.Gender
		data.Id = user.Id
		data.Nickname = user.Nickname
	} else {
		code = api.CodeUserExists
		err = errors.New("user already exists")
	}
	return
}

func (req *ReqLogIn) LogIn() (data RespLogin, code int, err error) {
	var user = repo.User{
		UserName: req.Username,
	}
	if user.FindByUsername() == nil {
		realPassword := security.DecryptWithPublicKey(user.Password, user.PublicKey)
		if realPassword == req.Password { //} security.DecryptWithPrivateKey(req.Password,user.PrivateKey) {
			data.Token, err = verify.SignToken(strconv.FormatInt(user.Id, 10))
			if err == nil {
				data.Id = user.Id
				data.Avatar = user.Avatar
				data.Username = user.UserName
				data.School = user.School
				data.StudentId = user.StudentId
				data.Gender = user.Gender
				data.Signature = user.Signature
				data.PublicKey = user.PublicKey
				data.Nickname = user.Nickname
			}
		} else {
			err = errors.New("wrong password")
			code = api.CodeWrongPassword
		}
	} else {
		err = errors.New("user does not exist")
		code = api.CodeUserNotExist
	}
	return
}
