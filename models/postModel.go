package models

import (
	"DoodleDropsBackend/traits"

	"gorm.io/gorm"
)

// TODO add image
type Post struct {
	gorm.Model
	ID          uint64 `gorm:"primaryKey"`
	AuthorID    uint64
	Author      User
	LikedBy     []*User `gorm:"many2many:user_liked_posts"`
	Description string
	Tags        []*Tag `gorm:"many2many:post_tags"`
	Comments    []*Comment
}

func (p *Post) PaginatedResult(db *gorm.DB, paginate *traits.Paginate) *gorm.DB {
	offset := (paginate.Page - 1) * paginate.Limit
	return db.Offset(offset).Limit(paginate.Limit)
}
