package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetHeaderUserId(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Keys["userId"].(string), 10, 64)
}
