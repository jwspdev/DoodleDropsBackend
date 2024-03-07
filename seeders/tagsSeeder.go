package seeders

import (
	"DoodleDropsBackend/models"

	"gorm.io/gorm"
)

//TODO MAKE THIS AS A SCRIPT THAT CAN BE RAN AS A COMMAND

func SeedTags(db *gorm.DB) error {
	tags := []*models.Tag{
		{Name: "Science", TagType: "tag", Description: "test"},
		{Name: "AI art", TagType: "tag", Description: "test"},
		{Name: "AI tool", TagType: "tool", Description: "test"},
		{Name: "MS Paint", TagType: "tool", Description: "test"},
		{Name: "Anime", TagType: "tag", Description: "test"},
		{Name: "Digital", TagType: "tag", Description: "test"},
		{Name: "Paper", TagType: "tool", Description: "test"},
		{Name: "Canvas", TagType: "tool", Description: "test"},
		{Name: "Abstract", TagType: "tag", Description: "test"},
		{Name: "Painting", TagType: "tag", Description: "test"},
		{Name: "Pens", TagType: "tool", Description: "test"},
		{Name: "Brush", TagType: "tool", Description: "test"},
		{Name: "Architecture", TagType: "tag", Description: "test"},
		{Name: "Impressionism", TagType: "tag", Description: "test"},
		{Name: "Pencil", TagType: "tool", Description: "test"},
		{Name: "Watercolor", TagType: "tool", Description: "test"},
	}
	for _, tag := range tags {
		if err := db.Create(tag).Error; err != nil {
			return err
		}
	}
	return nil
}
