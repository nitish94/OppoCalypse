package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth check for login page and static files
		if c.Request.URL.Path == "/login" || strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Next()
			return
		}

		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			fmt.Println("No user in session, redirecting to login")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Add user ID to context for use in handlers
		c.Set("userID", user)
		c.Next()
	}
}
