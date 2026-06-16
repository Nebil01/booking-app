package main

import (
	"fmt"
	"sync"
	"time"
)

var conferenceName = "Go Conference"

const conferenceTickets = 50

var remainingTickets uint = 50
var bookings = make([]UserData, 0)

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var wg = sync.WaitGroup{}

func main() {

	greetUsers()

	firstName, lastName, userTickets, email := getUserInput()
	isValidName, isValidEmail, isValidTicketsNumber := validateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidTicketsNumber {

		bookTicket(userTickets, firstName, lastName, email)

		wg.Add(1)
		go sendTicket(userTickets, firstName, lastName, email)

		firstNames := getFirstNames()
		fmt.Printf("The first name of bookings are: %v.\n", firstNames)

		if remainingTickets == 0 {
			fmt.Println("Our Conference is booked out.")
			//break
		}
	} else {
		if !isValidName {
			fmt.Println("Your first name or last name you entered is too short.")
		}
		if !isValidEmail {
			fmt.Println("Your email doesn't contain @ sign.")
		}
		if !isValidTicketsNumber {
			fmt.Println("Number of tickets you entered is invalid")
		}
	}
	wg.Wait()
}

func greetUsers() {
	fmt.Printf("Welcome to %v booking application.\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend.")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func getUserInput() (string, string, uint, string) {
	var firstName string
	var lastName string
	var userTickets uint
	var email string

	// Ask user their info

	fmt.Print("Enter your first name:\n")
	fmt.Scan(&firstName)

	fmt.Print("Enter your last name:\n")
	fmt.Scan(&lastName)

	fmt.Print("Enter number of tickets:\n")
	fmt.Scan(&userTickets)

	fmt.Print("Enter your email address:\n")
	fmt.Scan(&email)

	return firstName, lastName, userTickets, email
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets -= userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of all bookings are: %v\n", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v.\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v.\n", remainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(5 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("##################")
	fmt.Printf("Sending ticket:\n %v \nto email address %v.\n", ticket, email)
	fmt.Println("##################")
	wg.Done()
}
