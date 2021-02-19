package controller

import (
	_ "errors"
	"github.com/gin-gonic/gin"
	"hita/config"
	"hita/repository"
	"hita/service"
	"hita/utils/api"
	"hita/utils/logger"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

func UploadAvatar(c *gin.Context) {
	result := api.StdResp{}
	file, err := c.FormFile("upload")
	if err == nil {
		filename := c.Keys["userId"].(string) + path.Ext(file.Filename)
		fullPath := path.Join(logger.GetCurrentPath(), "..") + "/" + config.AvatarPath + filename
		_ = os.MkdirAll(path.Dir(fullPath), os.ModePerm)
		err = c.SaveUploadedFile(file, fullPath)
		if err == nil {
			idInt, _ := strconv.ParseInt(c.Keys["userId"].(string), 10, 64)
			user := repository.User{
				Id: idInt,
			}
			err = user.ChangeUserAvatar(filename)
			if err == nil {
				result.Code = api.CodeSuccess
				result.Message = "上传成功"
				result.Data = filename
			} else {
				_ = os.Remove(fullPath) //删除文件
				result.Code = api.CodeUserNotExist
				result.Message = "用户不存在"
				result.Data = gin.H{"error": err}
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

func GetAvatar(c *gin.Context) {
	result := api.StdResp{}
	id, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	user := repository.User{
		Id: id,
	}
	err := user.FindById()
	if err == nil {
		fullPath := path.Join(logger.GetCurrentPath(), "..") + "/" + config.AvatarPath + user.Avatar
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

func GetBasicProfile(c *gin.Context) {
	result := api.StdResp{}
	id, err := api.GetHeaderUserId(c)
	if err == nil {
		data, code, err := service.GetBasicProfile(id)
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
