package service

import (
	"fmt"
	repo "hita/repository"
	"hita/utils/api"
	"strconv"
	"time"
)

type RespTopic struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      int64     `json:"avatar"`
	ArticleNum  int       `json:"articleNum"`
	CreateTime  time.Time `json:"createTime"`
}

type GetTopicsReq struct {
	Mode     string `form:"mode" json:"mode"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	Extra    string `form:"extra" json:"extra"`
}

func (req *GetTopicsReq) GetTopics(userId int64) (result []RespTopic, code int, error error) {
	var topics []repo.Topic
	err := fmt.Errorf("")
	switch req.Mode {
	case "hot":
		{
			topics, err = repo.GetHotTopics()
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	case "search":
		{
			topics, err = repo.SearchTopics(req.Extra, req.PageSize, req.PageNum)
			if err != nil {
				return nil, api.CodeOtherError, err
			}
		}
	}
	for _, topic := range topics {
		resp := RespTopic{
			Id:          topic.Id,
			ArticleNum:  topic.ArticleNum,
			Name:        topic.Name,
			Avatar:      topic.Avatar,
			Description: topic.Description,
			CreateTime:  topic.CreateTime,
		}
		result = append(result, resp)
	}
	code = api.CodeSuccess
	return
}

type GetTopicReq struct {
	TopicId string `form:"topicId" json:"topicId"`
}

func (req *GetTopicReq) GetTopic() (result RespTopic, code int, er error) {
	topicIdInt, err := strconv.ParseInt(req.TopicId, 10, 64)
	if err != nil {
		return RespTopic{}, api.CodeWrongParam, err
	}
	topic := repo.Topic{
		Id: topicIdInt,
	}
	er = topic.Get()
	if er != nil {
		code = api.CodeOtherError
		return
	}
	result = RespTopic{
		Id:          topic.Id,
		ArticleNum:  topic.ArticleNum,
		Name:        topic.Name,
		Avatar:      topic.Avatar,
		Description: topic.Description,
		CreateTime:  topic.CreateTime,
	}
	code = api.CodeSuccess
	return
}
