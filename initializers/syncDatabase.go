package initializers

import (
	"gin/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
