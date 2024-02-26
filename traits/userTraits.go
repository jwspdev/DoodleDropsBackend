package traits

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckIfUserExists(c *gin.Context) (interface{}, bool) {
	user, exists := c.Get("user")
	return user, exists
}

func PromptUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
}
