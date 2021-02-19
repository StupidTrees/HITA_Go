package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
)

func SignUp(c *gin.Context) {
	var req service.ReqSignUp
	var result api.StdResp
	err := c.ShouldBind(&req)
	if err != nil {
		result.Code = -1
		result.Message = "request param error!:" + err.Error()
	} else {
		if req.Gender != "MALE" && req.Gender != "FEMALE" && req.Gender != "OTHER" {
			result.Code = api.CodeWrongParam
			result.Message = "wrong param！"
		} else {
			result.Data, result.Code, err = req.SignUp()
			if err != nil {
				result.Message = err.Error()
				result.Data = err
			} else {
				result.Code = api.CodeSuccess
				result.Message = "success!"
			}
		}
	}
	//响应给客户端
	c.JSON(200, result)
}

func LogIn(c *gin.Context) {
	var req service.ReqLogIn
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		if len(req.Username) == 0 {
			result.Code = api.CodeWrongParam
			result.Message = "username shouldn't be empty!"
		} else {
			result.Data, result.Code, err = req.LogIn()
			if err == nil {
				result.Code = api.CodeSuccess
				result.Message = "success!"
			} else {
				result.Data = err
				result.Message = "login failed"
			}
		}
	}
	//响应给客户端
	c.JSON(200, result)
}
