package main

import (
	"context"
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

func getClient(config *oauth2.Config) *http.Client {
    //tokFile := "token.json"
    //tok, err := tokenFromFile(tokFile)
    tok := getTokenFromWeb(config)
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

func main() {
    ctx := context.Background()
    creds, err := os.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read cred file: %v", err)
    }

    config, err := google.ConfigFromJSON(creds, calendar.CalendarReadonlyScope)
    if err != nil {
        log.Fatalf("Unable to parse cred file to config: %v", err)
    }

    client := getClient(config)

    srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        log.Fatalf("Unable to retrieve Calendar client: %v", err)
    }

    t := time.Now().Format(time.RFC3339)
    events, err := srv.Events.List("primary").ShowDeleted(false).
        SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
    if err != nil {
        log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
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
}

