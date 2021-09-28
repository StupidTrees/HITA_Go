package service

import (
	"hita/repository"
	"hita/utils/api"
	"strconv"
)

type FollowReq struct {
	FollowingId string `form:"followingId" json:"followingId"`
	Follow      bool   `form:"follow" json:"follow"`
}

type FollowResp struct {
	Follow bool `form:"follow" json:"follow"`
}

func (req *FollowReq) FollowOrUnFollow(userId int64) (data FollowResp, code int, error error) {
	aId, err := strconv.ParseInt(req.FollowingId, 10, 64)
	if err != nil {
		return FollowResp{}, api.CodeWrongParam, err
	}
	follow := repository.Follow{
		UserId:      userId,
		FollowingId: aId,
	}
	if req.Follow {
		error = follow.Create()
	} else {
		error = follow.Delete()
	}
	if error != nil {
		code = api.CodeOtherError
	} else {
		code = api.CodeSuccess
		data = FollowResp{
			Follow: req.Follow,
		}
	}
	return
}
