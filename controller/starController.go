package controller

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
)

func StarOrUnStar(c *gin.Context) {
	var req service.StarReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.StarOrUnStar(id)
		if err == nil {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		} else {
			result.Data = err
			result.Message = "create failed"

		}
	}
	fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}
