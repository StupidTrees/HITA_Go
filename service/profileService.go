package service

import (
	"errors"
	repo "hita/repository"
	"hita/utils/api"
	"strconv"
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
	Followed  bool   `json:"followed"`
}

type ProfileReq struct {
	UserId string `form:"userId" json:"userId"`
}

func (p *ProfileReq) GetBasicProfile(userId int64) (data RespUserProfile, code int, err error) {
	targetIdInt, e := strconv.ParseInt(p.UserId, 10, 64)
	if e != nil {
		return RespUserProfile{}, api.CodeWrongParam, e
	}
	var user = repo.User{Id: targetIdInt}
	var followed = false
	if targetIdInt != userId {
		follow := repo.Follow{
			UserId:      userId,
			FollowingId: targetIdInt,
		}
		followed = follow.Exist()
	}
	if user.FindById() == nil {
		data.Nickname = user.Nickname
		data.Id = user.Id
		data.Gender = user.Gender
		data.Username = user.UserName
		data.StudentId = user.StudentId
		data.Avatar = user.Avatar
		data.Signature = user.Signature
		data.School = user.School
		data.Followed = followed
		code = api.CodeSuccess
	} else {
		err = errors.New("user does not exist")
		code = api.CodeUserNotExist
	}
	return
}
