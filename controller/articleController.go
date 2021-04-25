package controller

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"hita/repository"
	"hita/service"
	"hita/utils/api"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

func CreateArticle(c *gin.Context) {
	var req service.CreateArticleReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Code, err = req.CreateArticle(id, []int64{})
		if err == nil {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		} else {
			result.Data = err
			result.Message = "create failed"
		}
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}

func CreateArticleWithImages(c *gin.Context) {
	result := api.StdResp{}
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["files"]
		var imageIds []int64
		for _, file := range files {
			filename := xid.New().String() + path.Ext(file.Filename)
			fullPath := repository.GetArticleImagePath(filename)
			_ = os.MkdirAll(path.Dir(fullPath), os.ModePerm)
			err = c.SaveUploadedFile(file, fullPath)
			if err == nil {
				img := repository.Image{
					Filename: filename,
					Type:     "POST",
				}
				err = img.Create()
				if err != nil {
					result.Code = -2
					result.Data = gin.H{"error": err}
					result.Message = "创建图片对象出错"
				} else {
					imageIds = append(imageIds, img.Id)
				}
			} else {
				result.Code = -2
				result.Data = gin.H{"error": err}
				result.Message = "上传文件出错"
			}
		}
		id, _ := api.GetHeaderUserId(c)
		req := service.CreateArticleReq{
			Content:  c.Query("content"),
			RepostId: c.Query("repostId"),
		}
		_, err = req.CreateArticle(id, imageIds)
		if err != nil {
			result.Code = api.CodeOtherError
			result.Data = gin.H{"error": err}
			result.Message = "发表文章失败"
		} else {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		}
	} else {
		result.Code = -1
		result.Data = gin.H{"error": err}
		result.Message = "接收表格出错"
	}
	//fmt.Println(result)
	c.JSON(http.StatusOK, result)
}

func GetArticles(c *gin.Context) {
	var req service.GetFollowingArticleReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetFollowingArticle(id)
		if err == nil {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		} else {
			result.Data = err
			result.Message = "create failed"

		}
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}

func GetArticle(c *gin.Context) {
	var req service.GetArticleReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetArticle(id)
		if err == nil {
			result.Code = api.CodeSuccess
			result.Message = "success!"
		} else {
			result.Data = err
			result.Message = "create failed"

		}
	}
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}
func GetArticleImage(c *gin.Context) {
	result := api.StdResp{}
	id, _ := strconv.ParseInt(c.Query("imageId"), 10, 64)
	image := repository.Image{
		Id: id,
	}
	err := image.Find()
	if err == nil {
		fullPath := repository.GetArticleImagePath(image.Filename)
		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Transfer-Encoding", "binary")
		data, err := ioutil.ReadFile(fullPath)
		if err == nil {
			c.Data(http.StatusOK, "image/jpeg", data)
		} else {
			result.Data = gin.H{"error": err}
			result.Message = "open file failed"
			result.Code = api.CodeOtherError
			c.JSON(http.StatusOK, result)
		}
	} else {
		result.Data = gin.H{"error": err}
		result.Message = "user not exist"
		result.Code = api.CodeUserNotExist
		c.JSON(http.StatusOK, result)
	}
}

func LikeOrUnlike(c *gin.Context) {
	var req service.LikeReq
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
	//fmt.Println(result)
	//响应给客户端
	c.JSON(200, result)
}

func DeleteArticle(c *gin.Context) {
	var req service.DeleteArticleReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Code, err = req.DeleteArticle(id)
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
