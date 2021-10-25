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
	Content    string `form:"content" json:"content"`
	TopicId    string `form:"topicId" json:"topicId"`
	RepostId   string `form:"repostId" json:"repostId"`
	AsAttitude bool   `form:"asAttitude" json:"asAttitude"`
	Anonymous  bool   `form:"anonymous" json:"anonymous"`
}

func (req *CreateArticleReq) CreateArticle(userId int64, imageIds []int64) (code int, error error) {
	article := repository.Article{
		Content:   req.Content,
		AuthorId:  userId,
		Images:    imageIds,
		Anonymous: req.Anonymous,
	}
	if req.AsAttitude {
		article.Type = "VOTE"
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
			if repost.AuthorId != userId{
				msg := repository.Message{
					UserId:      repost.AuthorId,
					Content:     article.Content,
					Type:        "ARTICLE",
					Action:      "REPOST",
					ReferenceId: req.RepostId,
					OtherId:     userId,
					CreateTime:  time.Now(),
				}
				_ = msg.Create()
			}
			if repost.RepostId > 0 {
				article.RepostId = repost.RepostId
			} else {
				article.RepostId = repost.Id
			}
		}
	}
	if len([]rune(req.TopicId)) > 0 {
		topicIdInt, err := strconv.ParseInt(req.TopicId, 10, 64)
		if err != nil {
			return api.CodeWrongParam, err
		}
		article.TopicId.Int64 = topicIdInt
		article.TopicId.Valid = true
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
	Id                 int64                `json:"id"`
	AuthorId           int64                `json:"authorId"`
	AuthorName         string               `json:"authorName"`
	AuthorAvatar       int64                `json:"authorAvatar"`
	RepostId           string               `json:"repostId"`
	RepostAuthorId     string               `json:"repostAuthorId"`
	RepostAuthorAvatar int64                `json:"repostAuthorAvatar"`
	RepostAuthorName   string               `json:"repostAuthorName"`
	RepostContent      string               `json:"repostContent"`
	RepostImages       repository.MIntArray `json:"repostImages"`
	RepostTime         time.Time            `json:"reposeTime"`
	TopicId            string               `json:"topicId"`
	TopicName          string               `json:"topicName"`
	Content            string               `json:"content"`
	Images             repository.MIntArray `json:"images"`
	Type               string               `json:"type"`
	Anonymous          bool                 `json:"anonymous"`
	LikeNum            int                  `json:"likeNum"`
	Liked              bool                 `json:"liked"`
	VotedUp            string               `json:"votedUp"`
	UpNum              int                  `json:"upNum"`
	IsMine             bool                 `json:"isMine"`
	DownNum            int                  `json:"downNum"`
	Starred            bool                 `json:"starred"`
	CommentNum         int                  `json:"commentNum"`
	CreateTime         time.Time            `json:"createTime"`
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

			articles, error = repository.SearchPosts(req.BeforeTime, req.AfterTime, req.PageSize, req.Extra)

		}
	case "user":
		{
			userIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			var tmpArticles []repository.Article
			tmpArticles, error = repository.GetUsersPosts(userIdInt, req.BeforeTime, req.AfterTime, req.PageSize)
			if userId==userIdInt{//看的不是自己的，把匿名隐藏
				articles = tmpArticles
			}else{
				for _,a := range tmpArticles {
					if !a.Anonymous{
						articles = append(articles, a)
					}
				}
			}
		}
	case "repost":
		{
			articleIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			articles, error = repository.GetReposts(articleIdInt, req.BeforeTime, req.AfterTime, req.PageSize)
		}
	case "star":
		{
			userIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			articles, error = repository.GetStarredPosts(userIdInt, req.BeforeTime, req.AfterTime, req.PageSize)
		}

	case "topic":
		{
			topicIdInt, err := strconv.ParseInt(req.Extra, 10, 64)
			if err != nil {
				return nil, api.CodeWrongParam, err
			}
			articles, error = repository.GetTopicPosts(topicIdInt, req.BeforeTime, req.AfterTime, req.PageSize)
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
			votedUp := "NONE"
			if a.Type == "VOTE" {
				vt := repository.Vote{
					UserId:    userId,
					ArticleId: a.Id,
				}
				err := vt.Find()
				if err == nil {
					if vt.Up {
						votedUp = "UP"
					} else {
						votedUp = "DOWN"
					}
				}
			}
			star := repository.Star{
				UserId:    userId,
				ArticleId: a.Id,
			}
			var topicId = ""
			if a.TopicId.Valid {
				topicId = fmt.Sprint(a.TopicId.Int64)
			}

			articleFormed := ArticleResp{
				Id:           a.Id,
				AuthorId:     a.Author.Id,
				AuthorName:   a.Author.Nickname,
				AuthorAvatar: a.Author.Avatar,
				Content:      a.Content,
				Images:       a.Images,
				TopicId:      topicId,
				TopicName:    a.Topic.Name,
				LikeNum:      a.LikeNum,
				UpNum:        a.UpNum,
				DownNum:      a.DownNum,
				Type:         a.Type,
				Liked:        like.Exist(),
				Anonymous:    a.Anonymous,
				VotedUp:      votedUp,
				Starred:      star.Exist(),
				RepostId:     "",
				IsMine:       userId == a.AuthorId,
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
					articleFormed.RepostAuthorId = fmt.Sprint(repostFrom.Author.Id)
					articleFormed.RepostTime = repostFrom.CreateTime
					articleFormed.RepostAuthorName = repostFrom.Author.Nickname
					articleFormed.RepostAuthorAvatar = repostFrom.Author.Avatar
					articleFormed.RepostImages = repostFrom.Images
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
	star := repository.Star{
		UserId:    userId,
		ArticleId: a.Id,
	}
	votedUp := "NONE"
	if a.Type == "VOTE" {
		vt := repository.Vote{
			UserId:    userId,
			ArticleId: a.Id,
		}
		err := vt.Find()
		if err == nil {
			if vt.Up {
				votedUp = "UP"
			} else {
				votedUp = "DOWN"
			}
		}
	}

	var topicId = ""
	if realObj.TopicId.Valid {
		topicId = fmt.Sprint(realObj.TopicId.Int64)
	}

	result = ArticleResp{
		Id:           realObj.Id,
		AuthorId:     realObj.Author.Id,
		AuthorName:   realObj.Author.Nickname,
		AuthorAvatar: realObj.Author.Avatar,
		Content:      realObj.Content,
		Images:       realObj.Images,
		LikeNum:      realObj.LikeNum,
		UpNum:        realObj.UpNum,
		DownNum:      realObj.DownNum,
		Type:         realObj.Type,
		TopicId:      topicId,
		Anonymous:    realObj.Anonymous,
		TopicName:    realObj.Topic.Name,
		Liked:        like.Exist(),
		VotedUp:      votedUp,
		IsMine:       userId == realObj.AuthorId,
		Starred:      star.Exist(),
		CommentNum:   realObj.CommentNum,
		CreateTime:   realObj.CreateTime,
	}
	if realObj.RepostId > 0 && !req.DigOrigin {
		repostFrom := repository.Article{
			Id: realObj.RepostId,
		}
		err := repostFrom.Get()
		if repostFrom.Anonymous {
			repostFrom.AuthorId = 0
		}
		if err == nil {
			result.RepostId = fmt.Sprint(repostFrom.Id)
			result.RepostContent = repostFrom.Content
			result.RepostAuthorId = fmt.Sprint(repostFrom.Author.Id)
			result.RepostTime = repostFrom.CreateTime
			result.RepostAuthorName = repostFrom.Author.Nickname
			result.RepostAuthorAvatar = repostFrom.Author.Avatar
			result.RepostImages = repostFrom.Images
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

	if error == nil && req.Like {
		a := repository.Article{
			Id: aId,
		}
		error = a.Get()
		if error == nil && a.AuthorId != userId {
			msg := repository.Message{
				UserId:      a.AuthorId,
				OtherId:     userId,
				Action:      "LIKE",
				Type:        "ARTICLE",
				ReferenceId: fmt.Sprint(a.Id),
				CreateTime:  userLike.CreateTime,
			}
			_ = msg.Create()
		}
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

type VoteReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
	Up        bool   `form:"up" json:"up"`
}

type VoteResp struct {
	UpNum   int    `form:"upNum" json:"upNum"`
	DownNum int    `form:"downNum" json:"downNum"`
	VotedUp string `form:"votedUp" json:"votedUp"`
}

func (req *VoteReq) Vote(userId int64) (data VoteResp, code int, error error) {
	aId, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return VoteResp{}, api.CodeWrongParam, err
	}
	userVote := repository.Vote{
		UserId:    userId,
		ArticleId: aId,
		Up:        req.Up,
	}
	error = userVote.Create()
	if error != nil {
		code = api.CodeOtherError
	} else {
		updated, err2 := userVote.GetVoteNum()
		votedUp := "UP"
		if !req.Up {
			votedUp = "DOWN"
		}
		if err2 == nil {
			code = api.CodeSuccess
			data = VoteResp{
				UpNum:   updated.UpNum,
				DownNum: updated.DownNum,
				VotedUp: votedUp,
			}
		} else {
			code = api.CodeOtherError
		}

	}
	return
}

type DeleteArticleReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
}

func (req *DeleteArticleReq) DeleteArticle(userId int64) (code int, error error) {
	idInt, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return api.CodeWrongParam, err
	}
	article := repository.Article{
		Id: idInt,
	}
	err = article.Get()
	if err != nil {
		return api.CodeArticleNotExist, err
	}
	if article.AuthorId != userId {
		return api.CodePermissionDenied, fmt.Errorf("不是你的帖子！")
	}
	//删除图片文件
	for _, id := range article.Images {
		img := repository.Image{
			Id: id,
		}
		_ = img.Delete()
	}
	err = article.Delete()
	if err != nil {
		return api.CodeOtherError, err
	}
	return api.CodeSuccess, nil
}
