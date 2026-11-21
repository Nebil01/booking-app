package handler

import (
	"booking-app/initializers"
	"booking-app/model"
	"booking-app/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BookTransaction(c *gin.Context) {
	var get struct {
		UserTicket uint `json:"user_ticket" binding:"required,gt=0"`
		EventID    uint `json:"event_id" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&get); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	userSession, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser, ok := userSession.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve user session"})
		return
	}

	var book model.Event

	booking := model.Work{
		Ticket:  get.UserTicket,
		UserID:  currentUser.ID,
		EventID: get.EventID,
	}
	// Basic transaction
	err := initializers.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, get.EventID).Error; err != nil {
			return err
		}

		if get.UserTicket > book.RemainingTickets {
			return errors.New("not enough tickets available")
		}

		book.RemainingTickets -= get.UserTicket
		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go service.SendTicket(currentUser.Email, book.Title, book.Date, get.UserTicket, booking.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Tickets booked successfully!"})
}

func GetTicket(c *gin.Context) {
	var ticket model.Work
	id := c.Param("id")
	rawUserID, _ := c.Get("user")
	user, ok := rawUserID.(*model.User)
	if !ok {
		c.JSON(401, gin.H{
			"error": "failed to retrieve user from context",
		})
		return
	}
	userID := user.ID
	if err := initializers.DB.Preload("Event").Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket": ticket,
	})
}

func GetMyTickets(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := rawUser.(*model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session format"})
		return
	}

	var tickets []model.Work
	if err := initializers.DB.Preload("Event").Where("user_id = ?", user.ID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

func GetAllEvents(c *gin.Context) {
	var events []model.Event
	if err := initializers.DB.Order("date ASC").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
