package controller

import (
	"time"

	"github.com/gin-gonic/gin"
)

func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}

type response struct {
	Message string `json:"message"`
}

func timeNow() *time.Time {
	now := time.Now()
	return &now
}
