package webhooks

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"sso-dashboard.bcgov.com/aggregator/utils"
)

type RocketChatNotifier interface {
	NotifyRocketChat(text string, title string, body string)
}

type RocketChat struct{}

/*
For the RC webhook, text will appear at the top, with a title and collapsible body below.
*/
func (r *RocketChat) NotifyRocketChat(text string, title string, body string) {
	// HTTP endpoint
	log.Println("Sending rocket chat notification")
	posturl := utils.GetEnv("RC_WEBHOOK", "")

	// JSON body
	requestBodyTemplate := `{
		"text": "%s",
		"attachments": [
			{
				"title": "%s",
				"text": "%s"
			}
		]
	}`

	requestBody := []byte(fmt.Sprintf(requestBodyTemplate, text, title, body))

	// Create a HTTP post request
	resp, err := http.NewRequest("POST", posturl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error sending rocket chat notification", err)
	}
	resp.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(resp)

	if err != nil {
		log.Println("Error sending rocket chat notification", err)
	}
	if res.StatusCode != 200 {
		log.Println("Error sending rocket chat notification", res.Status, res.Body)
	}

	defer res.Body.Close()
}
