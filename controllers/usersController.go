package controllers

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
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//Get the email and password
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	//Create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	userProfile := models.UserProfile{}

	userProfile.UserId = user.ID

	profileResult := initializers.DB.Create(&userProfile)

	if profileResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user ", "message": profileResult.Error})
		return
	}

	type userStruct struct {
		// Id    uint
		Email string
	}

	userResponse := userStruct{
		// Id:    user.ID,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    userResponse,
		"message": "User Created Successfully",
	})
}

func Login(c *gin.Context) {
	//Get email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	//Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or Password"})
		return
	}
	//Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or Password"})
		return
	}
	//Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Create Token"})
		return
	}
	//set cookie
	c.SetSameSite(http.SameSiteLaxMode)

	//send it back
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Logged In", "token": tokenString})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User is current logged in",
	})
}

func GetCurrentUser(c *gin.Context) {
	// Get the current logged-in user from the context
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, "User not found!")
		return
	}
	// Convert the interface to models.User
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, "Invalid User Type")
		return
	}
	if err := initializers.DB.Preload("UserProfile").Where("id = ?", currentUser.ID).First(&currentUser).Error; err != nil {
		fmt.Println("Error preloading user profile", err)
	}

	c.JSON(http.StatusOK, gin.H{"user": currentUser})
}
