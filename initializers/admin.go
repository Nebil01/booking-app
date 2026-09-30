package initializers

import (
	"booking-app/model"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	adminFName := os.Getenv("ADMIN_FIRST_NAME")
	adminLName := os.Getenv("ADMIN_LAST_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		return
	}

	var count int64
	DB.Model(&model.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), 10)
		adminUser := model.User{
			FirstName: adminFName,
			LastName:  adminLName,
			Email:     adminEmail,
			Password:  string(hashedPassword),
			Role:      "admin",
		}
		DB.Create(&adminUser)
		log.Println("✅ Initial admin created from .env")
	}
}
