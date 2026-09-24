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

func CreateEvent(c *gin.Context) {
	var input struct {
		Title         string `json:"title"`
		Date          string `json:"date"`
		TotalCapacity uint   `json:"total_capacity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	event := model.Event{
		Title:            input.Title,
		TotalCapacity:    input.TotalCapacity,
		RemainingTickets: input.TotalCapacity,
	}

	initializers.DB.Create(&event)
	c.JSON(200, gin.H{"message": "Event created successfully", "event": event})
}

func BookTransaction(c *gin.Context) {
	var get struct {
		UserTicket uint `json:"user_ticket"`
	}

	if err := c.ShouldBindJSON(&get); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Basic transaction
	err := initializers.DB.Transaction(func(tx *gorm.DB) error {
		var book model.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, 1).Error; err != nil {
			return err
		}

		if get.UserTicket > book.RemainingTickets {
			return errors.New("not enough tickets available")
		}

		book.RemainingTickets -= get.UserTicket
		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		userSession, _ := c.Get("user")
		currentUser := userSession.(*model.User)

		booking := model.Work{
			Ticket: get.UserTicket,
			UserID: currentUser.ID,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		go service.SendTicket(currentUser.Email, book.Title, book.Date, get.UserTicket)

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tickets booked successfully!"})
}

func GetEvent(c *gin.Context) {
	id := c.Param("id")
	var remaining model.Event

	if err := initializers.DB.First(&remaining, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"event_title":       remaining.Title,
		"total_capacity":    remaining.TotalCapacity,
		"remaining_tickets": remaining.RemainingTickets,
	})
}
