package middleware

import (
	"go-file/common"
	"go-file/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func permissionCheckHelper(c *gin.Context, requiredPermission int) {
	c.Set("role", common.RoleGuestUser)
	session := sessions.Default(c)
	role := session.Get("role")
	username := session.Get("username")
	if username != nil {
		c.Set("username", username)
	} else {
		// Check token
		token := c.Request.Header.Get("Authorization")
		user := model.ValidateUserToken(token)
		if user != nil && user.Username != "" {
			// Token is valid
			username = user.Username
			role = user.Role
			c.Set("username", username)
		} else {
			c.Set("username", "Guest user")
		}
	}
	if requiredPermission == common.RoleGuestUser {
		c.Next()
		return
	}
	if role == nil || role.(int) < requiredPermission {
		if c.Request.URL.Path == "/explorer" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"message":  "You do not have permission to access this page, please check whether you are logged in or have relevant permissions",
				"option":   common.OptionMap,
				"username": c.GetString("username"),
			})
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "You do not have permission to perform this operation. Please check whether you are logged in or have relevant permissions.",
			})
		}
		c.Abort()
		return
	}
	c.Set("role", role)
	c.Next()
}

func ImageDownloadPermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		permissionCheckHelper(c, common.ImageDownloadPermission)
	}
}

func ImageUploadPermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		permissionCheckHelper(c, common.ImageUploadPermission)
	}
}

func FileDownloadPermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		permissionCheckHelper(c, common.FileDownloadPermission)
	}
}

func FileUploadPermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		permissionCheckHelper(c, common.FileUploadPermission)
	}
}
