package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
)

func CountUserNum(c *gin.Context) {
	var result api.StdResp
	var err error
	result.Data, result.Code, err = service.CountUsers()
	if err == nil {
		result.Code = api.CodeSuccess
		result.Message = "success!"
	} else {
		result.Data = err
		result.Message = "create failed"
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}

func GetLatestVersionName(c *gin.Context) {
	var result api.StdResp
	var err error
	result.Data, result.Code, err = service.GetInfo("latest_version_name")
	if err == nil {
		result.Code = api.CodeSuccess
		result.Message = "success!"
	} else {
		result.Data = err
		result.Message = "fetch failed"
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}

func MakeSuggestion(c *gin.Context) {
	var req service.MakeSuggestionReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	result.Code, err = req.CreateSuggestion()
	if err == nil {
		result.Code = api.CodeSuccess
		result.Message = "success!"
	} else {
		result.Data = err
		result.Message = "fetch failed"
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}