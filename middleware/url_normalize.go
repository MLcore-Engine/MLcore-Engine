package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// URLNormalizer 处理URL末尾斜杠问题的中间件
func URLNormalizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 对API路径特殊处理
		if strings.HasPrefix(path, "/api/") && strings.HasSuffix(path, "/") {
			newPath := strings.TrimSuffix(path, "/")
			// 修改请求路径但不重定向
			c.Request.URL.Path = newPath
		}

		c.Next()
	}
}
