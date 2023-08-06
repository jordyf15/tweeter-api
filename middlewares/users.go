package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnsureCurrentUserIDMatchesPath(c *gin.Context) {
	if c.MustGet("current_user_id") != c.Param("user_id") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
