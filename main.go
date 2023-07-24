package main

import (
	"context"
    "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type application struct {
    cfg    *oauth2.Config
    client *http.Client
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

func createEvent(desc, start, end string) *calendar.Event {
    // Currently limited to: Description, Start and End time
    newEvent := calendar.Event{
        Description: desc, 
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

func (app *application) handleRedirectURI(w http.ResponseWriter, r *http.Request) {
    // TODO:
    // - make this handler only run the auth code stuff if it detects the "code" in the URL
    // - otherwise just serve the home page of this application

	// Get the "code" query parameter from the redirect URI
	code := r.URL.Query().Get("code")

	if code != "" {
		fmt.Println("OAuth2 Authorization Code:", code)
	} else {
		fmt.Println("No OAuth2 Authorization Code found in the redirect URI.")
	}

	// Handle other logic based on the authorization code if needed
	tok, err := app.cfg.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
    app.client = app.cfg.Client(context.Background(), tok)
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

    newEvent := createEvent("new event", t, time.Now().Add(time.Hour * 1).Format(time.RFC3339))
    evt, err := srv.Events.Insert("primary", newEvent).Do()
    if err != nil {
        log.Fatalf("Unable to create new event in calendar: %v", err)
    }
    fmt.Println("New event created: ")
	fmt.Printf("%v \n", evt.Description)
    }

func main() {
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
    mux.HandleFunc("/cal", app.calendarStuff)
    log.Print("Starting server on :4000")
    err = http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}

