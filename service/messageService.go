package service

import (
	"fmt"
	repo "hita/repository"
	"hita/utils/api"
	"strconv"
	"strings"
	"time"
)

type CountUnreadReq struct {
	Mode string `form:"mode" json:"mode"`
}

func (req *CountUnreadReq) CountUnread(userId int64) (result int64, code int, er error) {
	switch req.Mode {
	case "like":
		{
			result, er = repo.CountUnread(userId, "LIKE")
		}
	case "comment":
		{
			result, er = repo.CountUnread(userId, "COMMENT")
		}
	case "repost":
		{
			result, er = repo.CountUnread(userId, "REPOST")
		}
	case "follow":
		{
			result, er = repo.CountUnread(userId, "FOLLOW")
		}
	case "all":
		{
			result, er = repo.CountAllUnread(userId)
		}
	}
	if er != nil {
		code = api.CodeOtherError
		return
	}
	code = api.CodeSuccess
	return
}

type GetMessageReq struct {
	Mode     string `form:"mode" json:"mode"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
}

type MessageResp struct {
	Id            int64     `json:"id"`
	UserId        int64     `json:"userId"`
	UserName      string    `json:"userName"`
	UserAvatar    int64     `json:"userAvatar"`
	Action        string    `json:"action"`
	ActionContent string    `json:"actionContent"`
	Type          string    `json:"type"`
	ReferenceId   string    `json:"referenceId"`
	Content       string    `json:"content"`
	Image         string    `json:"Image"`
	CreateTime    time.Time `json:"createTime"`
}

func (req *GetMessageReq) GetMessages(userId int64) (result []MessageResp, code int, err error) {
	var messages []repo.Message
	switch req.Mode {
	case "like", "comment", "repost", "follow":
		{
			messages, err = repo.GetMessages(userId, strings.ToUpper(req.Mode), req.PageSize, req.PageNum)
		}
	}
	var ids []int64
	for _, msg := range messages {
		ids = append(ids, msg.Id)
		host := repo.User{
			Id: msg.OtherId,
		}
		err = host.FindById()
		if err != nil {
			continue
		}
		mr := MessageResp{
			Id:            msg.Id,
			UserId:        host.Id,
			UserName:      host.Nickname,
			UserAvatar:    host.Avatar,
			Action:        msg.Action,
			Type:          msg.Type,
			ReferenceId:   msg.ReferenceId,
			CreateTime:    msg.CreateTime,
			ActionContent: msg.Content,
		}
		switch msg.Type {
		case "COMMENT":
			{
				cmtIdInt, er := strconv.ParseInt(msg.ReferenceId, 10, 64)
				if er != nil {
					continue
				}
				cmt := repo.Comment{
					Id: cmtIdInt,
				}
				er = cmt.Get()
				if er != nil {
					continue
				}
				mr.Content = cmt.Content
			}
		case "ARTICLE":
			{
				articleIdInt, er := strconv.ParseInt(msg.ReferenceId, 10, 64)
				if er != nil {
					continue
				}
				a := repo.Article{
					Id: articleIdInt,
				}
				er = a.Get()
				if er != nil {
					continue
				}
				mr.Content = a.Content
				if len(a.Images) > 0 {
					mr.Image = fmt.Sprint(a.Images[0])
				}
			}
		}
		result = append(result, mr)
	}
	_ = repo.MarkAllRead(ids)


	code = api.CodeSuccess
	return
}
