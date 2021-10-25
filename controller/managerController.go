package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
)

func  CheckUpdate(c *gin.Context) {
	var req service.CheckUpdateReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	id, err := api.GetHeaderUserId(c)
	result.Data, result.Code, err = req.CheckUpdate(id)
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
