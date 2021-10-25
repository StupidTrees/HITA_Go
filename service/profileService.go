package service

import (
	"errors"
	"fmt"
	repo "hita/repository"
	"hita/utils/api"
	"strconv"
)

type RespUserProfile struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Gender       string `json:"gender"`
	Avatar       int64  `json:"avatar"`
	Signature    string `json:"signature"`
	StudentId    string `json:"studentId"`
	School       string `json:"school"`
	Followed     bool   `json:"followed"`
	FollowingNum int16  `json:"followingNum"`
	FansNum      int16  `json:"fansNum"`
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
		data.FansNum = user.FansNum
		data.FollowingNum = user.FollowingNum
		code = api.CodeSuccess
	} else {
		err = errors.New("user does not exist")
		code = api.CodeUserNotExist
	}
	return
}

type GetUserReq struct {
	Mode     string `form:"mode" json:"mode"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	Extra    string `form:"extra" json:"extra"`
}

func (req *GetUserReq) GetUser(userId int64) (result []RespUserProfile, code int, error error) {
	var users []repo.User
	switch req.Mode {
	case "liked":
		{
			articleIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			users, err = repo.GetLikedUsers(articleIdInt, req.PageSize, req.PageNum)
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	case "fans":
		{
			userIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			users, err = repo.GetFans(userIdInt, req.PageSize, req.PageNum)
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	case "following":
		{
			userIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			users, err = repo.GetFollowing(userIdInt, req.PageSize, req.PageNum)
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	case "search":
		{
			err := fmt.Errorf("")
			users, err = repo.SearchUser(req.Extra, req.PageSize, req.PageNum)
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	}
	for _, user := range users {
		var followed = false
		if user.Id != userId {
			follow := repo.Follow{
				UserId:      userId,
				FollowingId: user.Id,
			}
			followed = follow.Exist()
		}
		resp := RespUserProfile{
			Nickname:     user.Nickname,
			Id:           user.Id,
			Gender:       user.Gender,
			Username:     user.UserName,
			StudentId:    user.StudentId,
			Avatar:       user.Avatar,
			Signature:    user.Signature,
			School:       user.School,
			Followed:     followed,
			FollowingNum: user.FollowingNum,
			FansNum:      user.FansNum,
		}
		result = append(result, resp)
	}
	code = api.CodeSuccess
	return
}
