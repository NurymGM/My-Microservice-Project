package migrations

import (
	"github.com/NurymGM/jwtnew/initializers"
	"github.com/NurymGM/jwtnew/models"
)

func Migrate() {
	initializers.DB.AutoMigrate(&models.Userr{})
}
