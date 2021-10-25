package service

import (
	"hita/repository"
	"hita/utils/api"
	"strconv"
)

type StarReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
	Star      bool   `form:"star" json:"star"`
}

type StarResp struct {
	Starred bool `form:"starred" json:"starred"`
}

func (req *StarReq) StarOrUnStar(userId int64) (data StarResp, code int, error error) {
	aId, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return StarResp{}, api.CodeWrongParam, err
	}
	Star := repository.Star{
		UserId:    userId,
		ArticleId: aId,
	}
	if req.Star {
		error = Star.Create()
	} else {
		error = Star.Delete()
	}
	if error != nil {
		code = api.CodeOtherError
	} else {
		code = api.CodeSuccess
		data = StarResp{
			Starred: req.Star,
		}
	}
	return
}
