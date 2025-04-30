package webhooks

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"sso-dashboard.bcgov.com/aggregator/utils"
)

type RocketChatNotifier interface {
	NotifyRocketChat(text string, title string, body string) error
}

type RocketChat struct{}

/*
For the RC webhook, text will appear at the top, with a title and collapsible body below.
*/
func (r *RocketChat) NotifyRocketChat(text string, title string, body string) error {
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
	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error sending rocket chat notification", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("Error sending rocket chat notification", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Error sending rocket chat notification", res.Status, res.Body)
		return fmt.Errorf("received non-200 response: %s", res.Status)
	}

	defer res.Body.Close()
	return nil
}
