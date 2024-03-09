package controllers

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/requests"
	"DoodleDropsBackend/traits"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE AND DELETE COMMENT ONLY
func CreateComment(c *gin.Context) {
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

	request := requests.CreateCommentRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post = &models.Post{}
	if err := initializers.DB.First(&post, &request.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := &models.Comment{
		AuthorID: currentUser.ID,
		PostId:   &request.PostId,
		Content:  request.Content,
	}

	if err := initializers.DB.Create(comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "comment created successfully",
		"data":    comment,
	})
}

func DeleteComment(c *gin.Context) {
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

	var request struct {
		CommentId uint64 `json:"comment_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comment = &models.Comment{}
	if err := initializers.DB.First(&comment, request.CommentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not Found!"})
		return
	}
	if comment.AuthorID != currentUser.ID {
		traits.PromptUnauthorized(c, "Invalid User, cannot update")
		return
	}

	if err := initializers.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Post Deleted Successfully",
	})
}
