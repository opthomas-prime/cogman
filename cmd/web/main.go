package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type application struct {
	cfg    *oauth2.Config
	client *http.Client
}

func createEvent(sum, start, end string) *calendar.Event {
	// Currently limited to: Summary, Start and End time
	newEvent := calendar.Event{
		Summary: sum,
		Start: &calendar.EventDateTime{
			DateTime: start,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
		},
	}
	return &newEvent
}

// Basic Workflow for Data retrieval
// 1) Get all events from 00:00 to 23:59 and store as a 'Day' data structure linked to a user account
// 1a) Filter events from the past day
// 1b) Create users
// 1c) Create db for users


func main() {
	// Google creds 
	creds, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read cred file: %v", err)
	}

	config, err := google.ConfigFromJSON(creds, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse cred file to config: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	app := application{
		cfg: config,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleRedirectURI)
	mux.HandleFunc("/favicon.ico", doNothing)
	mux.HandleFunc("/cal", app.calendarStuff)
	log.Print("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
