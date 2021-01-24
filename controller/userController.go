package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
	"hita/utils/logger"
)

func SignUp(c *gin.Context) {
	var req service.ReqSignUp
	var token string
	var id int64
	resultCode := 0
	message := ""
	err := c.ShouldBind(&req)
	if err != nil {
		resultCode = -1
		message = "request param error!:" + err.Error()
		logger.Errorln(message)
	} else {
		if req.Gender != "MALE" && req.Gender != "FEMALE" && req.Gender != "OTHER" {
			resultCode = api.CodeWrongParam
			message = "wrong param！"
		} else {
			id, token, err = req.SignUp()
			if err != nil {
				resultCode = api.CodeUserExists
				message = "user already exists!"
			} else {
				message = "success!"
			}
		}
	}
	//响应给客户端
	c.JSON(200, api.StdResp{
		Code:    resultCode,
		Message: message,
		Data: gin.H{
			"token":  token,
			"userId": id,
		},
	})
}

func LogIn(c *gin.Context) {
	var req service.ReqLogIn
	var token string
	resultCode := 0
	message := ""
	err := c.ShouldBind(&req)
	if err != nil {
		resultCode = api.CodeWrongParam
		message = "request param error!:" + err.Error()
	} else {
		if len(req.Username) == 0 {
			resultCode = api.CodeWrongParam
			message = "username shouldn't be empty!"
		} else {
			token, resultCode, err = req.LogIn()
			if err == nil {
				resultCode = api.CodeSuccess
				message = "success!"
			} else {
				message = err.Error()
			}
		}
	}
	//响应给客户端
	c.JSON(200, api.StdResp{
		Code: resultCode, Message: message, Data: gin.H{
			"token": token,
		},
	})
}
