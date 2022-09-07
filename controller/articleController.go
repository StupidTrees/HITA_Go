package controller

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"hita/repository"
	"hita/service"
	"hita/utils/api"
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
		var images []repository.Image
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
				fmt.Println("检测敏感图...", repository.GetArticleImagePath(filename))
				res, _ := service.CheckSensitive(repository.GetArticleImagePath(filename))
				fmt.Println("检测结果：", res)
				img.Sensitive = res
				err = img.Create()
				if err != nil {
					result.Code = -2
					result.Data = gin.H{"error": err}
					result.Message = "创建图片对象出错"
				} else {
					imageIds = append(imageIds, img.Id)
					images = append(images, img)
				}
			} else {
				result.Code = -2
				result.Data = gin.H{"error": err}
				result.Message = "上传文件出错"
			}
		}
		id, _ := api.GetHeaderUserId(c)

		aa, _ := strconv.ParseBool(c.Query("asAttitude"))
		an, _ := strconv.ParseBool(c.Query("anonymous"))
		req := service.CreateArticleReq{
			Content:    c.Query("content"),
			RepostId:   c.Query("repostId"),
			TopicId:    c.Query("topicId"),
			Anonymous:  an,
			AsAttitude: aa,
		}
		_, err = req.CreateArticle(id, imageIds)
		if err != nil {
			result.Code = api.CodeOtherError
			result.Data = gin.H{"error": err}
			result.Message = "发表文章失败"
			for _, img := range images {
				_ = img.Delete()
			}
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
	data, err := service.GetImage(id)
	if err == nil {
		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Data(http.StatusOK, "image/jpeg", data)
	} else {
		result.Data = gin.H{"error": err}
		result.Message = "open file failed"
		result.Code = api.CodeOtherError
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

func Vote(c *gin.Context) {
	var req service.VoteReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.Vote(id)
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
