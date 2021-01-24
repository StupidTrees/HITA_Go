package api

import (
	"github.com/gin-gonic/gin"
)

type StdResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    gin.H  `json:"data"`
}