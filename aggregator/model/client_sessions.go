package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-co-op/gocron"
	"sso-dashboard.bcgov.com/aggregator/config"
	"sso-dashboard.bcgov.com/aggregator/keycloak"
	"sso-dashboard.bcgov.com/aggregator/utils"
	"sso-dashboard.bcgov.com/aggregator/webhooks"
)

type RealmInfo struct {
	Realm string `json:"realm"`
}

type SessionStats struct {
	ActiveSessions  string `json:"active"`
	ClientId        string `json:"clientId"`
	OfflineSessions string `json:"offline"`
}

var RealmErrorMessage = "Error getting realms for env %s: "
var ClientErrorMessage = "Error getting client stats for env %s: "

func GetRealms(rm *keycloak.RequestHandler) ([]string, error) {
	req, err := http.NewRequest("GET", rm.ApiBaseUrl+"/admin/realms", nil)
	if err != nil {
		log.Printf("Error occurred creating request: %v", err)
		return nil, err
	}

	body, err := rm.DoRequest(req)

	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	var realms []RealmInfo
	if err := json.Unmarshal([]byte(body), &realms); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return nil, err
	}

	var realmNames []string

	for _, realm := range realms {
		realmNames = append(realmNames, realm.Realm)
	}

	return realmNames, nil
}

func GetClientStats(rm *keycloak.RequestHandler, realms []string, env string) bool {
	errOccured := false

	handleError := func(message string) {
		log.Print(message)
		errOccured = true
	}

	for _, realm := range realms {
		req, err := http.NewRequest("GET", rm.ApiBaseUrl+"/admin/realms/"+realm+"/client-session-stats", nil)
		if err != nil {
			handleError(fmt.Sprintf("Error occurred creating request: %v", err))
			continue
		}

		body, err := rm.DoRequest(req)

		if err != nil {
			handleError(fmt.Sprintf("Error reading response body: %v", err))
			continue
		}

		var responseObjects []SessionStats
		err = json.Unmarshal(body, &responseObjects)
		if err != nil {
			handleError(fmt.Sprintf("Error unmarshaling JSON response: %v", err))
			continue
		}

		for _, obj := range responseObjects {
			err := InsertActiveSessions(env, realm, obj.ClientId, obj.ActiveSessions, obj.OfflineSessions)
			if err != nil {
				handleError(fmt.Sprintf("Error inserting active sessions: %v", err))
				continue
			}
		}
	}

	return errOccured
}

func InsertActiveSessions(environment string, realmID string, clientID string, activeSessions string, offlineSessions string) error {
	query := "INSERT INTO client_sessions (environment, realm_id, client_id, active_sessions, offline_sessions) VALUES(?,?,?,?,?);"
	_, err := pgdb.Query(nil, query, environment, realmID, clientID, activeSessions, offlineSessions)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ActiveSessions(env string, baseUrl string, clientId string, username string, password string, notifier webhooks.RocketChatNotifier) {
	log.Println("Getting active sessions for " + env + " environment")

	rm := keycloak.NewRequestHandler(&keycloak.RequestHandler{}, baseUrl, baseUrl, password, username, clientId)

	realms, err := GetRealms(rm)
	if err != nil {
		log.Println(fmt.Sprintf(RealmErrorMessage, env), err)
		notifier.NotifyRocketChat("Session Data Failure", fmt.Sprintf(RealmErrorMessage, env), err.Error())
		return
	}

	hasError := GetClientStats(rm, realms, strings.ToLower(env))
	if hasError {
		notifier.NotifyRocketChat("Session Data Failure", fmt.Sprintf(ClientErrorMessage, env), "One or more realm's client stats failed to be retrieved. See logs for details.")
		return
	}
	notifier.NotifyRocketChat("Session Data Loaded Successfully", env, fmt.Sprintf("Session data for environment %s has been loaded successfully", env))
}

func AllActiveSessions() {
	for _, env := range []string{"DEV", "TEST", "PROD"} {
		env := env
		baseUrl := utils.GetEnv(env+"_KEYCLOAK_URL", "")
		clientId := utils.GetEnv(env+"_KEYCLOAK_CLIENT_ID", "")
		username := utils.GetEnv(env+"_KEYCLOAK_USERNAME", "")
		password := utils.GetEnv(env+"_KEYCLOAK_PASSWORD", "")
		go func() {
			ActiveSessions(env, baseUrl, clientId, username, password, &webhooks.RocketChat{})
		}()
	}
}

func RunSessionsJob() {
	loc := config.LoadTimeLocation()
	cron := gocron.NewScheduler(loc)
	cron.Every(1).Hour().Do(AllActiveSessions)
	cron.StartAsync()
}
