package controllers

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/requests"
	"DoodleDropsBackend/responses"
	"DoodleDropsBackend/traits"
	"fmt"
	"net/http"

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

	post := &models.Post{AuthorID: currentUser.ID, Description: requestBody.Description, Title: requestBody.Title}
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
// TODO ADD MORE FIELDS
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
		PostId uint64 `json:"post_id"`
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
		CreatedAt:   post.CreatedAt,
	}
	c.JSON(http.StatusOK, response)

}

func ListPosts(c *gin.Context) {
	var posts []models.Post

	//parse query parameters for pagination
	// requestBody := requests.PaginateRequest{}
	var requestBody struct {
		UserID          *uint64                  `json:"user_id"`
		PaginateRequest requests.PaginateRequest `json:"paginate_request"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.PaginateRequest.Limit < 1 {
		requestBody.PaginateRequest.Limit = 5
	}

	if requestBody.PaginateRequest.PageNumber < 1 {
		requestBody.PaginateRequest.PageNumber = 1
	}

	paginate := &traits.Paginate{Limit: requestBody.PaginateRequest.Limit, Page: requestBody.PaginateRequest.PageNumber}

	//Query database for paginated posts
	db := initializers.DB.Model(&models.Post{})

	if requestBody.UserID != nil {
		// If UserID is provided, filter posts by author_id
		db = db.Where("author_id = ?", *requestBody.UserID)
	}

	if err := paginate.PagintedResult(db).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	var postResponses []responses.DefaultPostResponse
	for _, post := range posts {
		if err := initializers.DB.Preload("Author").Preload("Author.UserProfile").Preload("Tags").Preload("LikedBy").Preload("Comments").Omit("Author.Posts, Author.UserProfile.User").First(&post).Error; err != nil {
			fmt.Println("Error preloading post", err)
			return
		}
		UserLikedCount := len(post.LikedBy)
		CommentCount := len(post.Comments)
		var tagResponses []responses.TagsResponse
		for _, tag := range post.Tags {
			tagResponses = append(tagResponses, responses.TagsResponse{
				ID:      tag.ID,
				Name:    tag.Name,
				TagType: tag.TagType,
			})
		}

		postResponses = append(postResponses, responses.DefaultPostResponse{
			ID: post.ID,
			// CreatedAt:   post.CreatedAt,
			// UpdatedAt:   post.UpdatedAt,
			Title: post.Title,
			// Description: post.Description,
			Tags: tagResponses,
			PostCounters: responses.PostCounters{
				CommentCount:   CommentCount,
				UserLikedCount: UserLikedCount,
			},
			CreatedByResponse: responses.CreatedByResponse{
				UserID:      post.AuthorID,
				DisplayName: *post.Author.UserProfile.DisplayName,
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"page":       requestBody.PaginateRequest.PageNumber,
		"pageSize":   requestBody.PaginateRequest.Limit,
		"totalPosts": len(posts),
		"posts":      postResponses,
	})
}

// TODO FIX
// get current user's liked post paginated
// reason: only the current logged in user can see it's own liked posts, cannot see other user's liked post
func GetUsersLikedPost(c *gin.Context) {
	requestBody := requests.PaginateRequest{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, exists := traits.CheckIfUserExists(c)
	if requestBody.Limit < 1 {
		requestBody.Limit = 5
	}

	if requestBody.PageNumber < 1 {
		requestBody.PageNumber = 1
	}

	if !exists {
		traits.PromptUnauthorized(c, "User not found!")
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		traits.PromptUnauthorized(c, "Invalid User Type")
		return
	}
	if err := initializers.DB.Preload("LikedPosts").Find(&currentUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err})
		return
	}

	// paginate := &traits.Paginate{Limit: requestBody.Limit, Page: requestBody.PageNumber}
	likedPosts := currentUser.LikedPosts
	db := initializers.DB.Model(&likedPosts)
	c.JSON(http.StatusOK, gin.H{"posts": db.Find(&likedPosts)})
	// var paginatedLikedPosts []models.Post
	// if err := paginate.PagintedResult(db).Find(&paginatedLikedPosts).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
	// 	return
	// }

	// var postResponses []responses.DefaultPostResponse
	// for _, post := range likedPosts {
	// 	if err := initializers.DB.Preload("Author").Preload("Author.UserProfile").Preload("Tags").Preload("LikedBy").Preload("Comments").Omit("Author.Posts, Author.UserProfile.User").First(&post).Error; err != nil {
	// 		fmt.Println("Error preloading post", err)
	// 		return
	// 	}
	// 	UserLikedCount := len(post.LikedBy)
	// 	CommentCount := len(post.Comments)
	// 	var tagResponses []responses.TagsResponse
	// 	for _, tag := range post.Tags {
	// 		tagResponses = append(tagResponses, responses.TagsResponse{
	// 			ID:      tag.ID,
	// 			Name:    tag.Name,
	// 			TagType: tag.TagType,
	// 		})
	// 	}

	// 	postResponses = append(postResponses, responses.DefaultPostResponse{
	// 		ID:    post.ID,
	// 		Title: post.Title,
	// 		Tags:  tagResponses,
	// 		PostCounters: responses.PostCounters{
	// 			CommentCount:   CommentCount,
	// 			UserLikedCount: UserLikedCount,
	// 		},
	// 		CreatedByResponse: responses.CreatedByResponse{
	// 			UserID:      post.AuthorID,
	// 			DisplayName: *post.Author.UserProfile.DisplayName,
	// 		},
	// 	})
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"page":       requestBody.PageNumber,
	// 	"pageSize":   requestBody.Limit,
	// 	"totalPosts": len(currentUser.LikedPosts),
	// 	"posts":      postResponses,
	// 	// "posts": currentUser.LikedPosts,
	// })

}

//get paginated post user likers

func GetPostLikers(c *gin.Context) {
	var requestBody struct {
		PostID          uint64                   `json:"post_id"`
		PaginateRequest requests.PaginateRequest `json:"paginate_request"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestBody.PaginateRequest.Limit < 1 {
		requestBody.PaginateRequest.Limit = 5
	}

	if requestBody.PaginateRequest.PageNumber < 1 {
		requestBody.PaginateRequest.PageNumber = 1
	}
	post := &models.Post{}
	if err := initializers.DB.Preload("LikedBy").Omit("LikedBy").First(&post, requestBody.PostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting model"})
		return
	}
	paginate := &traits.Paginate{Limit: requestBody.PaginateRequest.Limit, Page: requestBody.PaginateRequest.PageNumber}
	db := initializers.DB.Model(&models.User{})

	if err := paginate.PagintedResult(db).Preload("UserProfile").Find(&post.LikedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	var likedByResponses []responses.LikedByResponse

	for _, user := range post.LikedBy {
		likedByResponses = append(likedByResponses, responses.LikedByResponse{
			DisplayName: *user.UserProfile.DisplayName,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"requestBody": likedByResponses,
	})
}
