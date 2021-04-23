package service

import (
	"fmt"
	"hita/repository"
	"hita/utils"
	"hita/utils/api"
	"strconv"
	"time"
)

type CreateArticleReq struct {
	Content  string `form:"content" json:"content"`
	RepostId string `form:"repostId" json:"repostId"`
}

func (req *CreateArticleReq) CreateArticle(userId int64) (code int, error error) {
	article := repository.Article{
		Content:  req.Content,
		AuthorId: userId,
	}
	if len([]rune(req.RepostId)) > 0 {
		repostIdInt, err := strconv.ParseInt(req.RepostId, 10, 64)
		if err != nil {
			return api.CodeWrongParam, err
		}
		repost := repository.Article{
			Id: repostIdInt,
		}
		err = repost.Get()
		if err == nil {
			if repost.RepostId > 0 {
				article.RepostId = repost.RepostId
			} else {
				article.RepostId = repost.Id
			}
		}
	}
	error = article.Create()
	if error == nil {
		code = api.CodeSuccess
	} else {
		code = api.CodeOtherError
	}
	return
}

type GetFollowingArticleReq struct {
	AfterTime  utils.Long `form:"afterTime" json:"afterTime"`
	BeforeTime utils.Long `form:"beforeTime" json:"beforeTime"`
	Mode       string     `form:"mode" json:"mode"`
	PageSize   int        `form:"pageSize" json:"pageSize"`
	Extra      string     `form:"extra" json:"extra"`
}
type ArticleResp struct {
	Id                 int64     `json:"id"`
	AuthorId           int64     `json:"authorId"`
	AuthorName         string    `json:"authorName"`
	AuthorAvatar       int64     `json:"authorAvatar"`
	RepostId           string    `json:"repostId"`
	RepostAuthorId     string    `json:"repostAuthorId"`
	RepostAuthorAvatar int64     `json:"repostAuthorAvatar"`
	RepostAuthorName   string    `json:"repostAuthorName"`
	RepostContent      string    `json:"repostContent"`
	RepostTime         time.Time `json:"reposeTime"`
	Content            string    `json:"content"`
	LikeNum            int       `json:"likeNum"`
	Liked              bool      `json:"liked"`
	CommentNum         int       `json:"commentNum"`
	CreateTime         time.Time `json:"createTime"`
}

func (req *GetFollowingArticleReq) GetFollowingArticle(userId int64) (result []ArticleResp, code int, error error) {

	var articles []repository.Article = nil
	switch req.Mode {
	case "following":
		{
			articles, error = repository.GetFollowingPosts(userId, req.BeforeTime, req.AfterTime, req.PageSize)
		}
	case "all":
		{
			articles, error = repository.GetAllPosts(req.BeforeTime, req.AfterTime, req.PageSize)
		}
	case "search":
		{
			{
				articles, error = repository.SearchPosts(req.BeforeTime, req.AfterTime, req.PageSize, req.Extra)
			}
		}
	}

	if error == nil {
		code = api.CodeSuccess
		var res []ArticleResp
		for _, a := range articles {
			like := repository.UserLikeArticle{
				UserId:    userId,
				ArticleId: a.Id,
			}
			articleFormed := ArticleResp{
				Id:           a.Id,
				AuthorId:     a.AuthorId,
				AuthorName:   a.Author.Nickname,
				AuthorAvatar: a.Author.Avatar,
				Content:      a.Content,
				LikeNum:      a.LikeNum,
				Liked:        like.Exist(),
				RepostId:     "",
				CommentNum:   a.CommentNum,
				CreateTime:   a.CreateTime,
			}
			if a.RepostId > 0 {
				repostFrom := repository.Article{
					Id: a.RepostId,
				}
				err := repostFrom.Get()
				if err == nil {
					articleFormed.RepostId = fmt.Sprint(repostFrom.Id)
					articleFormed.RepostContent = repostFrom.Content
					articleFormed.RepostAuthorId = fmt.Sprint(repostFrom.AuthorId)
					articleFormed.RepostTime = repostFrom.CreateTime
					articleFormed.RepostAuthorName = repostFrom.Author.Nickname
					articleFormed.RepostAuthorAvatar = repostFrom.Author.Avatar
				}
			}
			res = append(res, articleFormed)
		}
		result = res
	} else {
		code = api.CodeOtherError
	}
	return
}

type GetArticleReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
	DigOrigin bool   `form:"digOrigin" json:"digOrigin"`
}

func (req *GetArticleReq) GetArticle(userId int64) (result ArticleResp, code int, error error) {
	articleIdInt, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return ArticleResp{}, api.CodeWrongParam, err
	}
	code = api.CodeSuccess
	a := repository.Article{
		Id: articleIdInt,
	}
	err = a.Get()
	if err != nil {
		return ArticleResp{}, api.CodeOtherError, err
	}
	var realObj repository.Article
	if req.DigOrigin {
		realObj = repository.Article{
			Id: a.RepostId,
		}
		err = realObj.Get()
		if err != nil {
			return ArticleResp{}, api.CodeOtherError, err
		}
	} else {
		realObj = a
	}
	like := repository.UserLikeArticle{
		UserId:    userId,
		ArticleId: realObj.Id,
	}
	result = ArticleResp{
		Id:           realObj.Id,
		AuthorId:     realObj.AuthorId,
		AuthorName:   realObj.Author.Nickname,
		AuthorAvatar: realObj.Author.Avatar,
		Content:      realObj.Content,
		LikeNum:      realObj.LikeNum,
		Liked:        like.Exist(),
		CommentNum:   realObj.CommentNum,
		CreateTime:   realObj.CreateTime,
	}
	if realObj.RepostId > 0 && !req.DigOrigin {
		repostFrom := repository.Article{
			Id: realObj.RepostId,
		}
		err := repostFrom.Get()
		if err == nil {
			result.RepostId = fmt.Sprint(repostFrom.Id)
			result.RepostContent = repostFrom.Content
			result.RepostAuthorId = fmt.Sprint(repostFrom.AuthorId)
			result.RepostTime = repostFrom.CreateTime
			result.RepostAuthorName = repostFrom.Author.Nickname
			result.RepostAuthorAvatar = repostFrom.Author.Avatar
		}
	}
	return
}

type LikeReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
	Like      bool   `form:"like" json:"like"`
}

type LikeResp struct {
	LikeNum int  `form:"likeNum" json:"likeNum"`
	Liked   bool `form:"liked" json:"liked"`
}

func (req *LikeReq) LikeOrUnlike(userId int64) (data LikeResp, code int, error error) {
	aId, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return LikeResp{}, api.CodeWrongParam, err
	}
	userLike := repository.UserLikeArticle{
		UserId:    userId,
		ArticleId: aId,
	}
	if req.Like {
		error = userLike.Create()
	} else {
		error = userLike.Delete()
	}
	if error != nil {
		code = api.CodeOtherError
	} else {
		updated, err2 := userLike.GetLikeNum()
		if err2 == nil {
			code = api.CodeSuccess
			data = LikeResp{
				LikeNum: updated.LikeNum,
				Liked:   req.Like,
			}
		} else {
			code = api.CodeOtherError
		}

	}
	return
}
