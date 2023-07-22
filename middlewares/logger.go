package middlewares

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/color"
)

type LoggerMiddleware interface {
	PrintHeadersAndFormParams(c *gin.Context)
}

type loggerMiddleware struct {
}

func NewLoggerMiddleware() LoggerMiddleware {
	return &loggerMiddleware{}
}

func (middleware *loggerMiddleware) PrintHeadersAndFormParams(c *gin.Context) {
	if authorization := c.GetHeader("Authorization"); len(authorization) > 0 {
		fmt.Printf("%s[Headers]%s %s\n", color.Blue, color.Reset,
			color.Green+"Authorization: "+color.Reset+authorization)
	}

	c.PostForm("")
	for key, value := range c.Request.PostForm {
		if key == "password" {
			regex := regexp.MustCompile(`.`)
			value = []string{regex.ReplaceAllString(value[0], "*")}
		}

		fmt.Printf("%s[Form Data]%s %s%s:%s%s\n", color.Blue, color.Reset,
			color.Green, key, color.Reset, value)
	}
}
