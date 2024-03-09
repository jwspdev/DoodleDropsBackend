package models

import (
	"DoodleDropsBackend/traits"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"unique_index"`
	//limit tag type to tool or type
	TagType     string
	Description string
	LikedBy     []*User `gorm:"many2many:user_liked_tags"`
	//many to many
}

func (p *Tag) PaginatedResult(db *gorm.DB, paginate *traits.Paginate) *gorm.DB {
	offset := (paginate.Page - 1) * paginate.Limit
	return db.Offset(offset).Limit(paginate.Limit)
}
