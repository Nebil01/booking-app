package handler

import (
	"booking-app/initializers"
	"booking-app/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		FirstName string `json:"first_name" binding:"required,min=2"`
		LastName  string `json:"last_name" binding:"required,min=2"`
		Email     string `json:"email"`
		Password  string `json:"password" binding:"required,min=8"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash the password",
		})
		return
	}

	isValid := func(cond bool) bool {
		return cond
	}
	isValidEmail := isValid(strings.Contains(body.Email, "@") && strings.Contains(body.Email, "."))

	if !isValidEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your email doesn't contain @ or . sign.",
		})
		return
	}

	// create the user
	user := model.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hash),
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		message := result.Error.Error()
		if strings.Contains(strings.ToLower(message), "duplicate") || strings.Contains(strings.ToLower(message), "unique") {
			message = "Email already exists"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"Email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := strings.TrimSpace(os.Getenv("SECRET"))
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server misconfiguration: missing SECRET",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully!",
		"user":    tokenString,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user model.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID, //subject of the token, usually the user ID
		"Email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), //expire time
	})

	secret := strings.TrimSpace(os.Getenv("SECRET"))
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server misconfiguration: missing SECRET",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
