package controller

import (
	"encoding/json"
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/repository"
	"hita/service"
	"hita/utils/api"
	"strconv"
)

func Sync(c *gin.Context) {
	var req service.SyncReq
	var result api.StdResp
	err := c.ShouldBind(&req)
	if err != nil {
		result.Code = -1
		result.Message = "request param error!:" + err.Error()
	} else {
		result.Data, result.Code, err = req.Sync()
		if err != nil {
			result.Message = err.Error()
			result.Data = err
		} else {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		}
	}
	//fmt.Printf("%v\n", result.Data)
	//响应给客户端
	c.JSON(200, result)
}

func Push(c *gin.Context) {
	var req service.PushReq
	var result api.StdResp
	err := c.ShouldBind(&req)
	if err != nil {
		result.Code = -1
		result.Message = "request param error!:" + err.Error()
	} else {
		var historyList []repository.History
		_ = json.Unmarshal([]byte(req.History), &historyList)
		var dataMap map[string][]interface{}
		_ = json.Unmarshal([]byte(req.Data), &dataMap)
		uidI, _ := strconv.ParseInt(req.Uid, 10, 64)
		result.Code, err = req.Push(uidI, historyList, dataMap)
		if err != nil {
			result.Message = err.Error()
			result.Data = err
		} else {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		}
	}
	//响应给客户端
	c.JSON(200, result)
}
