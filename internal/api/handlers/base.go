package handlers

import (
	"github.com/gin-gonic/gin"
)

// Index 首页
func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
