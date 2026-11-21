package handler

import (
	"booking-app/initializers"
	"booking-app/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type adminUserResponse struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func toAdminUserResponse(user model.User) adminUserResponse {
	return adminUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ListUsers(c *gin.Context) {
	var users []model.User
	if err := initializers.DB.Select("id", "first_name", "last_name", "email", "role", "created_at").Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	response := make([]adminUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toAdminUserResponse(user))
	}
	c.JSON(http.StatusOK, gin.H{"users": response})
}

func CreateUser(c *gin.Context) {
	var input struct {
		FirstName string `json:"first_name" binding:"required,min=2,max=100"`
		LastName  string `json:"last_name" binding:"required,min=2,max=100"`
		Email     string `json:"email" binding:"required,email,max=254"`
		Password  string `json:"password" binding:"required,min=8,max=72"`
		Role      string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide valid name, email, password, and role values"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password could not be processed"})
		return
	}

	user := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      input.Role,
	}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not create user; the email may already be in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": toAdminUserResponse(user)})
}

func UpdateUser(c *gin.Context) {
	var input struct {
		FirstName string `json:"first_name" binding:"required,min=2,max=100"`
		LastName  string `json:"last_name" binding:"required,min=2,max=100"`
		Email     string `json:"email" binding:"required,email,max=254"`
		Password  string `json:"password" binding:"omitempty,min=8,max=72"`
		Role      string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide valid name, email, and role values"})
		return
	}

	var user model.User
	if err := initializers.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role == "admin" && input.Role != "admin" && isLastAdmin(initializers.DB) {
		c.JSON(http.StatusConflict, gin.H{"error": "The last admin cannot be demoted"})
		return
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.Role = input.Role
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password could not be processed"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not update user; the email may already be in use"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": toAdminUserResponse(user)})
}

func DeleteUser(c *gin.Context) {
	var user model.User
	if err := initializers.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	currentUser, ok := c.Get("user")
	admin, validAdmin := currentUser.(*model.User)
	if !ok || !validAdmin || admin.ID == user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own admin account"})
		return
	}
	if user.Role == "admin" && isLastAdmin(initializers.DB) {
		c.JSON(http.StatusConflict, gin.H{"error": "The last admin cannot be deleted"})
		return
	}

	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}

func isLastAdmin(db *gorm.DB) bool {
	var count int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)
	return count <= 1
}
