package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func (app *application) handleRedirectURI(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - make this handler only run the auth code stuff if it detects the "code" in the URL
	// - otherwise just serve the home page of this application

	// Get the "code" query parameter from the redirect URI
	code := r.URL.Query().Get("code")

	if code != "" {
		fmt.Println("OAuth2 Authorization Code:", code)
		// Handle other logic based on the authorization code if needed
		tok, err := app.cfg.Exchange(context.TODO(), code)
		if err != nil {
			w.Write([]byte("ERROR"))
			log.Fatalf("Unable to retrieve token from web: %v", err)
		}
		app.client = app.cfg.Client(context.Background(), tok)
	} else {
		fmt.Println("No OAuth2 Authorization Code found in the redirect URI.")
	}

	w.Write([]byte("OK"))
}

func (app *application) calendarStuff(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client := app.client
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	currentT := time.Now()
	t := currentT.Format(time.RFC3339)
	startOfDay := time.Date(currentT.Year(), currentT.Month(), currentT.Day(), 0, 0, 0, 0, currentT.Location())
	endOfDay := startOfDay.Add(time.Hour * 24).Add(-time.Second)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(startOfDay.Format(time.RFC3339)).
		TimeMax(endOfDay.Format(time.RFC3339)).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve past day of user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}

	newEvent := createEvent("New Event", t, time.Now().Add(time.Hour*1).Format(time.RFC3339))
	evt, err := srv.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Fatalf("Unable to create new event in calendar: %v", err)
	}
	fmt.Println("New event created: ")
	fmt.Printf("%v \n", evt.Summary)
}

func doNothing(w http.ResponseWriter, r *http.Request) {
}
