package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/lib/logger"
	"hita/service"
)


func SignUp(c *gin.Context){
	var req service.ReqSignUp
	var token int64
	resultCode := 0
	message := ""
	err := c.ShouldBind(&req)
	if err != nil {
		resultCode = -1
		message = "request param error!:" + err.Error()
		logger.Errorln(message)
	} else {
		if req.Gender != "MALE"&& req.Gender != "FEMALE" && req.Gender != "OTHER" {
			resultCode = CodeWrongParam
			message = "wrong param！"
		}else{
			token,err = req.SignUp()
			if err!=nil {
				resultCode = CodeUserExists
				message = "user already exists!"
			}
		}
	}
	//响应给客户端
	c.JSON(200,gin.H{
		"result_code": resultCode,
		"error_msg":   message,
		"data": token,
	})
}


