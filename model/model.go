package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID    uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string `json:"role" gorm:"default:'user'"`
}

type Work struct {
	gorm.Model
	Ticket  uint
	UserID  uint
	EventID uint
	Event   Event `json:"event"`
}

type Event struct {
	gorm.Model
	Title            string
	Date             string
	TotalCapacity    uint
	RemainingTickets uint
}
