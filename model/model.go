package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID    uint
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Work struct {
	gorm.Model
	Ticket   uint
	TicketID uint
}

type Event struct {
	gorm.Model
	Title            string
	Date             string
	TotalCapacity    uint
	RemainingTickets uint
}
