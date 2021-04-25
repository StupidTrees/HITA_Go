package controller

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hita/service"
	"hita/utils/api"
)

func CreateComment(c *gin.Context) {
	var req service.CreateCommentReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Code, err = req.CreateComment(id)
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
func GetCommentsOfArticle(c *gin.Context) {
	var req service.GetCommentsOfArticleReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetCommentsOfArticle(id)
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

func GetCommentsOfComment(c *gin.Context) {
	var req service.GetCommentsOfCommentReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetCommentsOfComment(id)
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

func GetCommentInfo(c *gin.Context) {
	var req service.GetCommentReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetCommentInfo(id)
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

func LikeOrUnlikeComment(c *gin.Context) {
	var req service.LikeCommentReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.LikeOrUnlike(id)
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

func DeleteComment(c *gin.Context) {
	var req service.DeleteCommentReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Code, err = req.DeleteComment(id)
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
