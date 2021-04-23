package service

import (
	"encoding/json"
	"fmt"
	"hita/repository"
	"hita/utils/api"
	"strconv"
	"time"
)

type CreateCommentReq struct {
	ArticleId  string `form:"articleId",json:"articleId"`
	ReceiverId string `form:"receiverId",json:"receiverId"`
	ContextId  string `form:"contextId",json:"contextId"`
	ReplyId    string `form:"replyId",json:"replyId"`
	Content    string `form:"content" json:"content"`
}

func (req *CreateCommentReq) CreateComment(userId int64) (code int, error error) {
	j, _ := json.Marshal(*req)
	fmt.Println(string(j))
	articleIdInt, er := strconv.ParseInt(req.ArticleId, 10, 64)
	if er != nil {
		return api.CodeWrongParam, er
	}

	recIdInt, er := strconv.ParseInt(req.ReceiverId, 10, 64)
	if er != nil {
		return api.CodeWrongParam, er
	}
	var repIdInt int64
	if len([]rune(req.ReplyId)) > 0 {
		repIdInt, error = strconv.ParseInt(req.ReplyId, 10, 64)
		if error != nil {
			return api.CodeWrongParam, error
		}
	}

	var contextIdInt int64
	if len([]rune(req.ContextId)) > 0 {
		contextIdInt, error = strconv.ParseInt(req.ContextId, 10, 64)
		if error != nil {
			return api.CodeWrongParam, error
		}
	}

	comment := repository.Comment{
		Content:    req.Content,
		ArticleId:  articleIdInt,
		AuthorId:   userId,
		ContextId:  contextIdInt,
		ReceiverId: recIdInt,
		ReplyId:    repIdInt,
	}
	error = comment.Create()
	if error == nil {
		code = api.CodeSuccess
	} else {
		code = api.CodeOtherError
	}
	return
}

type GetCommentsOfArticleReq struct {
	ArticleId string `form:"articleId" json:"articleId"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
}
type CommentResp struct {
	Id             int64     `json:"id"`
	ArticleId      int64     `json:"articleId"`
	AuthorId       string    `json:"authorId"`
	AuthorName     string    `json:"authorName"`
	AuthorAvatar   int64     `json:"authorAvatar"`
	ReceiverId     string    `json:"receiverId"`
	ReceiverName   string    `json:"receiverName"`
	ReceiverAvatar int64     `json:"receiverAvatar"`
	ReplyId        string    `json:"replyId"`
	ReplyContent   string    `json:"replyContent"`
	Content        string    `json:"content"`
	LikeNum        int       `json:"likeNum"`
	CommentNum     int       `json:"commentNum"`
	Liked          bool      `json:"liked"`
	CreateTime     time.Time `json:"createTime"`
}

func (req *GetCommentsOfArticleReq) GetCommentsOfArticle(userId int64) (result []CommentResp, code int, error error) {
	articleIdInt, err := strconv.ParseInt(req.ArticleId, 10, 64)
	if err != nil {
		return nil, api.CodeWrongParam, err
	}
	articles, err := repository.GetCommentsOfArticle(articleIdInt, req.PageSize, req.PageNum)
	if err == nil {
		code = api.CodeSuccess
		var res []CommentResp

		for _, a := range articles {
			var replyContent string
			if a.ReplyTo != nil {
				replyContent = a.ReplyTo.Content
			}
			var replyIdStr string
			if a.ReplyId <= 0 {
				replyIdStr = ""
			} else {
				replyIdStr = fmt.Sprint(a.ReplyId)
			}
			liked := repository.UserLikeComment{
				UserId:    userId,
				CommentId: a.Id,
			}
			res = append(res, CommentResp{
				Id:             a.Id,
				ArticleId:      a.ArticleId,
				AuthorId:       fmt.Sprint(a.AuthorId),
				AuthorName:     a.Author.Nickname,
				AuthorAvatar:   a.Author.Avatar,
				ReceiverId:     fmt.Sprint(a.ReceiverId),
				ReceiverName:   a.Receiver.Nickname,
				ReceiverAvatar: a.Receiver.Avatar,
				ReplyId:        replyIdStr,
				ReplyContent:   replyContent,
				Content:        a.Content,
				LikeNum:        a.LikeNum,
				CommentNum:     a.CommentNum,
				Liked:          liked.Exist(),
				CreateTime:     a.CreateTime,
			})
		}
		result = res
	} else {
		error = err
		code = api.CodeOtherError
	}
	return
}

type GetCommentsOfCommentReq struct {
	CommentId string `form:"commentId" json:"commentId"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
}

func (req *GetCommentsOfCommentReq) GetCommentsOfComment(userId int64) (result []CommentResp, code int, error error) {
	commentIdInt, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return nil, api.CodeWrongParam, err
	}
	articles, err := repository.GetCommentsOfComment(commentIdInt, req.PageSize, req.PageNum)
	if err == nil {
		code = api.CodeSuccess
		var res []CommentResp

		for _, a := range articles {
			var replyContent string
			if a.ReplyTo != nil {
				replyContent = a.ReplyTo.Content
			}
			var replyIdStr string
			if a.ReplyId <= 0 {
				replyIdStr = ""
			} else {
				replyIdStr = fmt.Sprint(a.ReplyId)
			}
			liked := repository.UserLikeComment{
				UserId:    userId,
				CommentId: a.Id,
			}
			res = append(res, CommentResp{
				Id:             a.Id,
				ArticleId:      a.ArticleId,
				AuthorId:       fmt.Sprint(a.AuthorId),
				AuthorName:     a.Author.Nickname,
				AuthorAvatar:   a.Author.Avatar,
				ReceiverId:     fmt.Sprint(a.ReceiverId),
				ReceiverName:   a.Receiver.Nickname,
				ReceiverAvatar: a.Receiver.Avatar,
				ReplyId:        replyIdStr,
				ReplyContent:   replyContent,
				Content:        a.Content,
				LikeNum:        a.LikeNum,
				CommentNum:     a.CommentNum,
				Liked:          liked.Exist(),
				CreateTime:     a.CreateTime,
			})
		}
		result = res
	} else {
		error = err
		code = api.CodeOtherError
	}
	return
}

type GetCommentReq struct {
	CommentId string `form:"commentId" json:"commentId"`
}

func (req *GetCommentReq) GetCommentInfo(userId int64) (result CommentResp, code int, error error) {
	commentIdInt, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return CommentResp{}, api.CodeWrongParam, err
	}
	a, err := repository.GetCommentInfo(commentIdInt)
	if err == nil {
		code = api.CodeSuccess
		var replyContent string
		if a.ReplyTo != nil {
			replyContent = a.ReplyTo.Content
		}
		var replyIdStr string
		if a.ReplyId <= 0 {
			replyIdStr = ""
		} else {
			replyIdStr = fmt.Sprint(a.ReplyId)
		}
		liked := repository.UserLikeComment{
			UserId:    userId,
			CommentId: a.Id,
		}
		result = CommentResp{
			Id:             a.Id,
			ArticleId:      a.ArticleId,
			AuthorId:       fmt.Sprint(a.AuthorId),
			AuthorName:     a.Author.Nickname,
			AuthorAvatar:   a.Author.Avatar,
			ReceiverId:     fmt.Sprint(a.ReceiverId),
			ReceiverName:   a.Receiver.Nickname,
			ReceiverAvatar: a.Receiver.Avatar,
			ReplyId:        replyIdStr,
			ReplyContent:   replyContent,
			Content:        a.Content,
			LikeNum:        a.LikeNum,
			Liked:          liked.Exist(),
			CreateTime:     a.CreateTime,
		}
	} else {
		error = err
		code = api.CodeOtherError
	}
	return
}

type LikeCommentReq struct {
	CommentId string `form:"commentId" json:"commentId"`
	Like      bool   `form:"like" json:"like"`
}

func (req *LikeCommentReq) LikeOrUnlike(userId int64) (data LikeResp, code int, error error) {
	aId, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return LikeResp{}, api.CodeWrongParam, err
	}
	userLike := repository.UserLikeComment{
		UserId:    userId,
		CommentId: aId,
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
