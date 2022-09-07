package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{})
}

func UserAgreementPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user_agreement.html", gin.H{})
}

func PrivacyPolicyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "privacy_policy.html", gin.H{})
}
