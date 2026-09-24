package service

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendTicket(toEmail string, eventTitle string, date string, ticketCount uint) error {
	from := mail.NewEmail("Project Hermes", "nebulibrahim@gmail.com")
	subject := "Your Ticket Confirmation for " + eventTitle
	to := mail.NewEmail("Recipient", toEmail)
	plainTextContent := "You have successfully booked " + fmt.Sprint(ticketCount) + " tickets for " + eventTitle + "."
	htmlContent := "<strong>You have successfully booked " + fmt.Sprint(ticketCount) + " tickets for " + eventTitle + "!</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error response: %s", response.Body)
		return errors.New("sendgrid returned error status: " + response.Body)
	}

	log.Printf("Email sent successfully! Status: %d", response.StatusCode)
	return nil
}
