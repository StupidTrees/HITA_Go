package controller

import (
	_ "errors"
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

func UploadAvatar(c *gin.Context) {
	result := api.StdResp{}
	file, err := c.FormFile("upload")
	if err == nil {
		filename := xid.New().String() + path.Ext(file.Filename)
		fullPath := repository.GetAvatarPath(filename)
		_ = os.MkdirAll(path.Dir(fullPath), os.ModePerm)
		err = c.SaveUploadedFile(file, fullPath)
		if err == nil {
			idInt, _ := strconv.ParseInt(c.Keys["userId"].(string), 10, 64)
			img := repository.Image{
				Filename: filename,
				Type:     "AVATAR",
			}
			err = img.Create()
			if err == nil {
				user := repository.User{
					Id: idInt,
				}
				err = user.FindById()
				if err == nil {
					oldImg := repository.Image{
						Id: user.Avatar,
					}
					_ = oldImg.Delete()
					err = user.ChangeUserAvatar(img.Id)
					if err == nil {
						result.Code = api.CodeSuccess
						result.Message = "上传成功"
						result.Data = img.Id
					} else {
						_ = img.Delete()
						result.Code = api.CodeUserNotExist
						result.Message = "用户不存在"
						result.Data = gin.H{"error": err}
					}
				} else {
					_ = img.Delete()
					result.Code = api.CodeUserNotExist
					result.Message = "用户不存在"
					result.Data = gin.H{"error": err}
				}

			} else {
				result.Code = -2
				result.Data = gin.H{"error": err}
				result.Message = "创建图片对象出错"
			}

		} else {
			result.Code = -2
			result.Data = gin.H{"error": err}
			result.Message = "上传文件出错"
		}
	} else {
		result.Code = -1
		result.Data = gin.H{"error": err}
		result.Message = "接收表格出错"
	}
	c.JSON(http.StatusOK, result)
}

func changeProfile(c *gin.Context, param string, attr string) {
	result := api.StdResp{}
	id, err := api.GetHeaderUserId(c)
	if err == nil {
		user := repository.User{Id: id}
		err = user.ChangeUserProfile(attr, c.PostForm(param))
		if err == nil {
			result.Message = "success"
			result.Code = api.CodeSuccess
		} else {
			result.Data = err
			result.Code = api.CodeOtherError
			result.Message = "operation failed"
		}
	} else {
		result.Data = err
		result.Code = api.CodeWrongParam
		result.Message = "wrong header"
	}
	c.JSON(http.StatusOK, result)
}

func ChangeSignature(c *gin.Context) {
	changeProfile(c, "signature", "signature")
}

func ChangeNickname(c *gin.Context) {
	changeProfile(c, "nickname", "nickname")
}

func ChangeGender(c *gin.Context) {
	var param = c.PostForm("gender")
	if param == "MALE" || param == "FEMALE" || param == "OTHER" {
		changeProfile(c, "gender", "gender")
	} else {
		c.JSON(http.StatusOK, api.StdResp{Code: api.CodeWrongParam, Message: "wrong param!"})
	}
}

func GetAvatar(c *gin.Context) {
	result := api.StdResp{}
	id, _ := strconv.ParseInt(c.Query("imageId"), 10, 64)
	data, err := service.GetImage(id)
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Transfer-Encoding", "binary")
	if err == nil {
		c.Data(http.StatusOK, "image/jpeg", data)
	} else {
		result.Data = gin.H{"error": err}
		result.Message = "open file failed or user not found"
		result.Code = api.CodeOtherError
		c.JSON(http.StatusOK, result)
	}
}

func GetBasicProfile(c *gin.Context) {
	req := service.ProfileReq{}
	result := api.StdResp{}
	err := c.ShouldBind(&req)
	if err == nil {
		id, err := api.GetHeaderUserId(c)
		data, code, err := req.GetBasicProfile(id)
		if err == nil {
			result.Data = data
			result.Message = "success"
			result.Code = api.CodeSuccess
		} else {
			result.Data = err
			result.Message = err.Error()
			result.Code = code
		}
	} else {
		result.Data = err
		result.Message = "wrong id format!"
		result.Code = api.CodeWrongParam
	}
	c.JSON(http.StatusOK, result)
}

func GetUsers(c *gin.Context) {
	var req service.GetUserReq
	var result api.StdResp
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		result.Code = api.CodeWrongParam
		result.Message = "request param error!"
	} else {
		id, err := api.GetHeaderUserId(c)
		result.Data, result.Code, err = req.GetUser(id)
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
