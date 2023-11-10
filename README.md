# cogman

## Notes

* Build: `go build` in the cmd/web directory
* Config: Go to
* https://console.cloud.google.com/apis/credentials?project=<YOUR_PROJECT>
    * You will need to create a project in GCP and replace <YOUR_PROJECT> with your project name
* Enable the Google Calendar API: https://console.cloud.google.com/apis/library/calendar-json.googleapis.com?project=<YOUR_PROJECT>
* Configure OAuth Consent Screen
    * I needed to add an authorized domain; I'm not completely sure if this is required though.
    * I needed to add a test user (myemail@gmail.com)
* Configure a OAuth 2.0 Client Id
    * Redirect uri: http://localhost:4000
    * Authorized JavaScript origins: http://localhost
    * Once you create the app, download the json of the config and move it to your current directory as `credentials.json`
* Run the app: `./web`
* Follow the link it gives
* If you get "OK" try going to: `http://localhost:4000/cal`
    * As of 11.09.23, this will print out your current day and add a test calendar item to your calendar
* If you want to revoke OAuth Access, go here: `https://myaccount.google.com/`
