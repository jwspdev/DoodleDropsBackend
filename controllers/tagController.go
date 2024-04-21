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

// listing and relating of tags
func ListTags(c *gin.Context) {
	var tags []models.Tag

	requestBody := requests.PaginateRequest{}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Limit < 1 {
		requestBody.Limit = 5
	}

	if requestBody.PageNumber < 1 {
		requestBody.PageNumber = 1
	}

	paginate := &traits.Paginate{Limit: requestBody.Limit, Page: requestBody.PageNumber}

	db := initializers.DB.Model(&models.Tag{})
	if err := paginate.PagintedResult(db).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	var tagsResponse []responses.TagsResponse
	for _, tag := range tags {
		tagsResponse = append(tagsResponse, responses.TagsResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			TagType:     tag.TagType,
		})
	}

	response := struct {
		Page      int                      `json:"page"`
		PageSize  int                      `json:"page_size"`
		TotalTags int                      `json:"total_tags"`
		Tags      []responses.TagsResponse `json:"tags"`
	}{
		Page:      requestBody.PageNumber,
		PageSize:  requestBody.Limit,
		TotalTags: len(tags),
		Tags:      tagsResponse,
	}

	c.JSON(http.StatusOK, response)
}

func AddTagToPost(c *gin.Context) {
	var request struct {
		PostId uint64   `json:"post_id"`
		TagIDs []uint64 `json:"tag_ids"`
	}
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
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var post = &models.Post{}
	if err := initializers.DB.First(&post, request.PostId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot find post"})
		return
	}
	if post.AuthorID != currentUser.ID {
		traits.PromptUnauthorized(c, "Invalid user, cannot execute!")
		return
	}
	var tags []models.Tag
	if err := initializers.DB.Find(&tags, request.TagIDs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tags Not Found"})
		return
	}
	var newTagIDs []uint64
	for _, selected := range tags {
		found := false
		for _, tag := range post.Tags {
			if tag.ID == selected.ID {
				found = true
				break
			}
		}
		if !found {
			newTagIDs = append(newTagIDs, selected.ID)
		}
	}

	db := initializers.DB

	tx := db.Begin()

	insertedTags := 0

	for _, tag := range tags {
		alreadyAdded := false
		for _, addedTag := range post.Tags {
			if addedTag.ID == tag.ID {
				alreadyAdded = true
				break
			}
		}

		if !alreadyAdded {
			query := `
                INSERT INTO post_tags (post_id, tag_id)
                SELECT ?, ?
                WHERE NOT EXISTS (
                    SELECT 1 FROM post_tags WHERE post_id = ? AND tag_id = ?
                )
            `
			result := tx.Exec(query, post.ID, tag.ID, post.ID, tag.ID)
			if result.Error != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			if result.RowsAffected > 0 {
				insertedTags++
			}
		}
	}
	tx.Commit()
	if insertedTags > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tags added successfully",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Tags Added",
		})
	}
}

func LikeTag(c *gin.Context) {
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
		TagIDs []uint64 `json:"tag_ids"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tags []models.Tag
	if err := initializers.DB.Find(&tags, request.TagIDs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tags Not Found"})
		return
	}

	var newTagIDs []uint64
	for _, selected := range tags {
		found := false
		for _, tag := range currentUser.LikedTags {
			if tag.ID == selected.ID {
				found = true
				break
			}
		}
		if !found {
			newTagIDs = append(newTagIDs, selected.ID)
		}
	}

	db := initializers.DB

	tx := db.Begin()

	insertedTags := 0

	for _, tag := range tags {
		alreadyLiked := false
		for _, likedTag := range currentUser.LikedTags {
			if likedTag.ID == tag.ID {
				alreadyLiked = true
				break
			}
		}

		if !alreadyLiked {
			query := `
                INSERT INTO user_liked_tags (user_id, tag_id)
                SELECT ?, ?
                WHERE NOT EXISTS (
                    SELECT 1 FROM user_liked_tags WHERE user_id = ? AND tag_id = ?
                )
            `
			result := tx.Exec(query, currentUser.ID, tag.ID, currentUser.ID, tag.ID)
			if result.Error != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			if result.RowsAffected > 0 {
				insertedTags++
			}
		}
	}

	// Commit the transaction
	tx.Commit()
	if insertedTags > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tags liked successfully",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Tags Added",
		})
	}

}
