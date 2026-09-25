package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendTicket(toEmail string, eventTitle string, date string, ticketCount uint, bookingID uint) error {
	from := mail.NewEmail("Project Hermes", "nebulibrahim@gmail.com")
	subject := "Your Ticket Confirmation for " + eventTitle
	to := mail.NewEmail("Recipient", toEmail)

	// 1. Generate QR Code file path and data string
	qrPath := fmt.Sprintf("./qrcodes/booking_%d.png", bookingID)
	qrData := fmt.Sprintf("Booking ID: %d | Event: %s | Date: %s | Tickets: %d", bookingID, eventTitle, date, ticketCount)

	// 2. Call GenerateQRCode
	if err := GenerateQRCode(qrData, qrPath); err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		return err
	}
	defer os.Remove(qrPath) // Clean up temp file after email sends

	// 3. Build Email Content
	plainTextContent := fmt.Sprintf("Thank you for your purchase! Booking ID: %d for %s.", bookingID, eventTitle)
	htmlContent := fmt.Sprintf("<strong>Thank you for your purchase!</strong><p>Event: %s<br>Date: %s<br>Tickets: %d</p><p>Your QR code ticket is attached below.</p>", eventTitle, date, ticketCount)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// 4. Read generated QR image and attach to SendGrid email
	fileBytes, err := os.ReadFile(qrPath)
	if err == nil {
		encoded := base64.StdEncoding.EncodeToString(fileBytes)
		attachment := mail.NewAttachment()
		attachment.SetContent(encoded)
		attachment.SetType("image/png")
		attachment.SetFilename(fmt.Sprintf("ticket_%d.png", bookingID))
		attachment.SetDisposition("attachment")
		message.AddAttachment(attachment)
	}

	// 5. Send Email via SendGrid
	client := sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error response: %s", response.Body)
		return errors.New("sendgrid returned error status: " + response.Body)
	}

	log.Printf("Email with QR code sent successfully! Status: %d", response.StatusCode)
	return nil
}
