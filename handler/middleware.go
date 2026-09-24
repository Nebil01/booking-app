package handler

import (
	"booking-app/initializers"
	"booking-app/model"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
		tokenString = strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	}

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		secret := os.Getenv("SECRET")
		if secret == "" {
			return nil, errors.New("missing secret")
		}
		return []byte(secret), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(expFloat) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	subFloat, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user model.User
	userID := int(subFloat)
	initializers.DB.First(&user, userID)
	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", &user)
	c.Set("userID", userID)
	c.Next()
}
