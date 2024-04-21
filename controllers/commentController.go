package controllers

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/models"
	"DoodleDropsBackend/requests"
	"DoodleDropsBackend/responses"
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

func GetPostComments(c *gin.Context) {
	var comments []models.Comment
	var requestBody struct {
		PostID          *uint64                  `json:"post_id"`
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

	db := initializers.DB.Model(&models.Comment{})
	db.Preload("Author.UserProfile").Omit("Author")
	db.Where("post_id = ?", &requestBody.PostID)

	if err := paginate.PagintedResult(db).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	var commentResponses []responses.CommentResponse

	for _, comment := range comments {
		commentResponses = append(commentResponses, responses.CommentResponse{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			Content:   comment.Content,
			CreatedBy: responses.CreatedByResponse{
				UserID:      comment.AuthorID,
				DisplayName: *comment.Author.UserProfile.DisplayName,
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"page":           requestBody.PaginateRequest.PageNumber,
		"page_size":      requestBody.PaginateRequest.Limit,
		"total_comments": len(comments),
		"comments":       commentResponses,
	})
}
