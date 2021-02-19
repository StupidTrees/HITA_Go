package service

import (
	"errors"
	repo "hita/repository"
	"hita/utils/api"
)

type RespUserProfile struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
	StudentId string `json:"studentId"`
	School    string `json:"school"`
}

func GetBasicProfile(userId int64) (data RespUserProfile, code int, err error) {
	var user = repo.User{Id: userId}
	if user.FindById() == nil {
		data.Nickname = user.Nickname
		data.Id = user.Id
		data.Gender = user.Gender
		data.Username = user.UserName
		data.StudentId = user.StudentId
		data.Avatar = user.Avatar
		data.Signature = user.Signature
		data.School = user.School
		code = api.CodeSuccess
	} else {
		err = errors.New("user does not exist")
		code = api.CodeUserNotExist
	}
	return
}
