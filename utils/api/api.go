package api

import (
	"github.com/gin-gonic/gin"
)

type StdResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}
