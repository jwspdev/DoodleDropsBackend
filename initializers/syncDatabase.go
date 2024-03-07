package initializers

import (
	"DoodleDropsBackend/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{}, &models.Comment{}, &models.UserProfile{})
	// DB.AutoMigrate(&models.Post{})
	// DB.AutoMigrate(&models.Comment{})

	DB.AutoMigrate(&models.PostTag{})

	DB.AutoMigrate(&models.UserLikedPost{})
	DB.AutoMigrate(&models.UserLikedTags{})

	// DB.Migrator().CreateConstraint(&models.User{}, "Posts")
	// DB.Migrator().CreateConstraint(&models.User{}, "fk_UserID")

	//seed
	// if err := seeders.SeedTags(DB); err != nil {
	// 	panic("Failed to seed to db")
	// }
	// DB.Migrator().CreateConstraint(&models.)
	// DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.UserProfile{})

	//MANUALLY ALTER TABLES (USED FOR TESTING IF MIGRATION IS PROPER)
	// DB.Exec("ALTER TABLE users DROP COLUMN profile_refer").Exec("ALTER TABLE users DROP COLUMN test")
}
