package service

import (
	repo "hita/repository"
	"hita/utils/api"
)




func  CountUsers() (result int64, code int, er error) {
	result,er = repo.CountUserNum()
	if er != nil {
		code = api.CodeOtherError
	}else{
		code = api.CodeSuccess
	}
	return
}

func  GetInfo(key string) (result string, code int, er error) {
	info:=repo.Info{
		Key:key,
	}
	er = info.Get()
	if er != nil {
		code = api.CodeOtherError
	}else{
		result = info.Value
		code = api.CodeSuccess
	}
	return
}


type MakeSuggestionReq struct {
	Content       string     `form:"content" json:"content"`
}
func (req *MakeSuggestionReq)CreateSuggestion() (code int, er error) {
	info:=repo.Inbox{
		Content:req.Content,
	}
	er = info.CreateSuggestion()
	if er != nil {
		code = api.CodeOtherError
	}else{
		code = api.CodeSuccess
	}
	return
}

