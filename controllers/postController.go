package controllers

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/requests"
	"DoodleDropsBackend/responses"
	"DoodleDropsBackend/traits"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const userNotFoundMessage string = "User not found!"
const invalidUserTypeMessage string = "Invalid User Type"

// Create
func CreatePost(c *gin.Context) {
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, userNotFoundMessage)
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, invalidUserTypeMessage)
		return
	}

	requestBody := requests.CreatePostRequest{}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &models.Post{AuthorID: currentUser.ID, Description: requestBody.Description}
	if err := initializers.DB.Create(post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	if err := initializers.DB.Preload("Author").First(&post).Error; err != nil {
		fmt.Println("Error preloading user profile", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post created successfully",
		"data":    post,
	})
}

// TODO CREATE A METHOD THAT GETS THE CURRENT POST AND CHECKS IF THE USER IS VALIDATED TO RETRIEVE SUCH POST
// Update
func UpdatePost(c *gin.Context) {
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, userNotFoundMessage)
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, invalidUserTypeMessage)
		return
	}

	var postIdRequest struct {
		PostId            uint64
		PostUpdateRequest requests.CreatePostRequest
	}
	if err := c.BindJSON(&postIdRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post = &models.Post{}
	if err := initializers.DB.First(&post, postIdRequest.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != currentUser.ID {
		traits.PromptUnauthorized(c, "Invalid user, cannot update")
		return
	}
	post.Description = postIdRequest.PostUpdateRequest.Description

	if err := initializers.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post Updated Successfully",
	})
}

// Delete
func DeletePost(c *gin.Context) {
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, userNotFoundMessage)
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, invalidUserTypeMessage)
		return
	}

	var postIdRequest struct {
		PostId uint64
	}
	if err := c.BindJSON(&postIdRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post = &models.Post{}
	if err := initializers.DB.First(&post, postIdRequest.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != currentUser.ID {
		traits.PromptUnauthorized(c, "Invalid user, cannot update")
		return
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post Deleted Successfully",
	})
}

//List (implement pagination)

// Get Current
func GetCurrentPost(c *gin.Context) {
	user, exists := traits.CheckIfUserExists(c)

	if !exists {
		traits.PromptUnauthorized(c, userNotFoundMessage)
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, invalidUserTypeMessage)
		return
	}

	var postIdRequest struct {
		PostId uint64
	}
	if err := c.BindJSON(&postIdRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var post = &models.Post{}
	if err := initializers.DB.First(&post, postIdRequest.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := initializers.DB.
		Preload("Author.UserProfile").
		Select("id, author_id, description").
		Omit("Author.Posts, Author.UserProfile.User").
		First(&post, postIdRequest.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	var isEditable bool = post.AuthorID == currentUser.ID
	response := responses.CurrentPostResponse{
		IsEditable:  isEditable,
		Description: post.Description,
		ID:          post.ID,
		CreatedBy: responses.CreatedByResponse{
			UserID:      post.Author.ID,
			DisplayName: *post.Author.UserProfile.DisplayName,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

func ListPosts(c *gin.Context) {
	var posts []models.Post

	//parse query parameters for pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "2"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if err != nil || limit < 1 {
		limit = 3
	}

	paginate := &traits.Paginate{Limit: limit, Page: page}

	//Query database for paginated posts
	db := initializers.DB.Model(&models.Post{})
	if err := paginate.PagintedResult(db).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":       page,
		"pageSize":   limit,
		"totalPosts": len(posts),
		"posts":      posts,
	})
}
