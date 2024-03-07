package controllers

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/responses"
	"DoodleDropsBackend/traits"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// [x] UPDATE USER PROFILE
// [x] ADD CHECKER IF USER HAS PROFILE, IF NOT, THEN CREATE ONE
// [x] IF USER PASSES A VALUE, THEN UPDATE IT, IF THEY DON'T, THEN RETAIN THE PREVIOUS RECORD.
func UpdateUserProfile(c *gin.Context) {
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, "User not found!")
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, "Invalid User Type")
		return
	}
	if err := initializers.DB.Preload("UserProfile").Where("id = ?", currentUser.ID).First(&currentUser).Error; err != nil {
		fmt.Println("Error preloading user profile", err)
	}

	var requestBody struct {
		DisplayName *string
		FirstName   *string
		MiddleName  *string
		LastName    *string
		Age         *uint8
		Birthday    *time.Time
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userProfile := currentUser.UserProfile
	if err := initializers.DB.Preload("User").Where("id = ?", userProfile.ID).First(&userProfile).Error; err != nil {
		fmt.Println("Error preloading user profile", err)
	}
	if userProfile == nil {
		userProfile = &models.UserProfile{UserId: currentUser.ID}
		err := initializers.DB.Create(userProfile).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err,
			})
		}
	}

	if requestBody.DisplayName != nil {
		userProfile.DisplayName = requestBody.DisplayName
	}

	if requestBody.FirstName != nil {
		userProfile.FirstName = requestBody.FirstName
	}

	if requestBody.MiddleName != nil {
		userProfile.MiddleName = requestBody.MiddleName
	}

	if requestBody.LastName != nil {
		userProfile.LastName = requestBody.LastName
	}

	if requestBody.Age != nil {
		userProfile.Age = requestBody.Age
	}

	if requestBody.Birthday != nil {
		userProfile.Birthday = requestBody.Birthday
	}

	err := initializers.DB.Save(userProfile).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
	}
	userResponse := responses.UserProfileResponse{
		UpdatedAt:   currentUser.UserProfile.UpdatedAt,
		DisplayName: currentUser.UserProfile.DisplayName,
		FirstName:   currentUser.UserProfile.FirstName,
		MiddleName:  currentUser.UserProfile.MiddleName,
		LastName:    currentUser.UserProfile.LastName,
		Age:         currentUser.UserProfile.Age,
		Birthday:    currentUser.UserProfile.Birthday,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"data":    userResponse,
	})
}
