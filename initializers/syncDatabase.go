package initializers

import "DoodleDropsBackend/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.UserProfile{})

	DB.Exec("ALTER TABLE users DROP COLUMN profile_refer").Exec("ALTER TABLE users DROP COLUMN test")
}
