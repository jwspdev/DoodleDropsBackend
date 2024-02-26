package middleware

import (
	// "DoodleDropsBackend/initializers"
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/traits"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	//Get the cookie off req
	// tokenString, err := c.Cookie("Authorization")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]
	//Decode/Validate it
	// Parse takes the token string and a function for looking up the key.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		// c.AbortWithStatus(http.StatusUnauthorized)
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		traits.PromptUnauthorized(c, "Invalid token")
		c.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}
		//Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
			c.Abort()
			return
		}
		//Attach to req
		c.Set("user", user)

		//continue
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		c.Abort()
	}
}
