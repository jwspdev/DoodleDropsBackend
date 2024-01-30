package initializers

import "DoodleDropsBackend/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
