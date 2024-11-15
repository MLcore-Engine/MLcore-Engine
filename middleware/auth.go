package middleware

import (
	"MLcore-Engine/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(minRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		tokenString := authHeader[7:] // Remove "Bearer " prefix
		claims, err := common.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		if claims.Role < minRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "无权进行此操作，权限不足",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// func authHelper(c *gin.Context, minRole int) {
// 	session := sessions.Default(c)
// 	username := session.Get("username")
// 	role := session.Get("role")
// 	id := session.Get("id")
// 	status := session.Get("status")
// 	authByToken := false
// 	if username == nil {
// 		// Check token
// 		token := c.Request.Header.Get("Authorization")
// 		if token == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"success": false,
// 				"message": "无权进行此操作，未登录或 token 无效",
// 			})
// 			c.Abort()
// 			return
// 		}
// 		user, err := model.ValidateUserToken(token)
// 		common.SysLog(err.Error())
// 		if user != nil && user.Username != "" {
// 			// Token is valid
// 			username = user.Username
// 			role = user.Role
// 			id = user.Id
// 			status = user.Status
// 		} else {
// 			c.JSON(http.StatusOK, gin.H{
// 				"success": false,
// 				"message": "无权进行此操作，token 无效",
// 			})
// 			c.Abort()
// 			return
// 		}
// 		authByToken = true
// 	}
// 	if status.(int) == common.UserStatusDisabled {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": false,
// 			"message": "用户已被封禁",
// 		})
// 		c.Abort()
// 		return
// 	}
// 	if role.(int) < minRole {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": false,
// 			"message": "无权进行此操作，权限不足",
// 		})
// 		c.Abort()
// 		return
// 	}
// 	c.Set("username", username)
// 	c.Set("role", role)
// 	c.Set("id", id)
// 	c.Set("authByToken", authByToken)
// 	c.Next()
// }

func UserAuth() func(c *gin.Context) {
	return JWTAuthMiddleware(common.RoleCommonUser)
}

func AdminAuth() func(c *gin.Context) {
	return JWTAuthMiddleware(common.RoleAdminUser)
}

func RootAuth() func(c *gin.Context) {
	return JWTAuthMiddleware(common.RoleRootUser)
}
