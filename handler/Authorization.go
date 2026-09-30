package handler

import (
	"booking-app/initializers"
	"booking-app/model"
	"net/http"

	"github.com/gin-gonic/gin"
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
		Date:             input.Date,
		TotalCapacity:    input.TotalCapacity,
		RemainingTickets: input.TotalCapacity,
	}

	if err := initializers.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event created successfully", "event": event})
}

func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var event model.Event

	if err := initializers.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var input struct {
		Title         string `json:"title"`
		Date          string `json:"date"`
		TotalCapacity uint   `json:"total_capacity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update payload"})
		return
	}

	initializers.DB.Model(&event).Updates(model.Event{
		Title:         input.Title,
		Date:          input.Date,
		TotalCapacity: input.TotalCapacity,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully", "event": event})
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&model.Event{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
