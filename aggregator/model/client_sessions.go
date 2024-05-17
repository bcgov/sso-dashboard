package model

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"io"

	"github.com/go-co-op/gocron"
	"sso-dashboard.bcgov.com/aggregator/config"
	"sso-dashboard.bcgov.com/aggregator/keycloak"
	"sso-dashboard.bcgov.com/aggregator/utils"
)

type RealmInfo struct {
	Realm string `json:"realm"`
}

type SessionStats struct {
	ActiveSessions  string `json:"active"`
	ClientId        string `json:"clientId"`
	OfflineSessions string `json:"offline"`
}

func GetRealms(tm *keycloak.TokenManager) []string {
	req, err := http.NewRequest("GET", tm.BaseUrl+"/admin/realms", nil)
	if err != nil {
		log.Fatalf("Error occurred creating request: %v", err)
	}

	resp, _ := tm.DoRequest(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var realms []RealmInfo
	if err := json.Unmarshal([]byte(body), &realms); err != nil {
		log.Fatalf("Error unmarshalling response body: %v", err)
	}

	var realmNames []string

	for _, realm := range realms {
		realmNames = append(realmNames, realm.Realm)
	}

	return realmNames
}

func GetClientStats(tm *keycloak.TokenManager, realms []string, env string) error {
	for _, realm := range realms {
		req, err := http.NewRequest("GET", tm.BaseUrl+"/admin/realms/"+realm+"/client-session-stats", nil)
		if err != nil {
			log.Fatalf("Error occurred creating request: %v", err)
			return nil
		}

		resp, _ := tm.DoRequest(req)

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
			return err
		}

		var responseObjects []SessionStats
		err = json.Unmarshal(body, &responseObjects)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON response: %v", err)
			return err
		}

		for _, obj := range responseObjects {
			err := InsertActiveSessions(env, realm, obj.ClientId, obj.ActiveSessions, obj.OfflineSessions)
			if err != nil {
				log.Fatalf("Error inserting active sessions: %v", err)
			}
		}
	}
	return nil
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

func ActiveSessions(env string) {
	log.Println("Getting active sessions for " + env + " environment")
	baseUrl := utils.GetEnv(env+"_KEYCLOAK_URL", "")
	clientId := utils.GetEnv(env+"_KEYCLOAK_CLIENT_ID", "")
	username := utils.GetEnv(env+"_KEYCLOAK_USERNAME", "")
	password := utils.GetEnv(env+"_KEYCLOAK_PASSWORD", "")

	tm := keycloak.NewTokenManager(clientId, password, username, baseUrl)

	realms := GetRealms(tm)
	GetClientStats(tm, realms, strings.ToLower(env))
}

func AllActiveSessions() {
	go func() {
		ActiveSessions("DEV")
	}()

	go func() {
		ActiveSessions("TEST")
	}()

	go func() {
		ActiveSessions("PROD")
	}()
}

func RunSessionsJob() {
	loc := config.LoadTimeLocation()
	cron := gocron.NewScheduler(loc)
	cron.Every(1).Hour().Do(AllActiveSessions)
	cron.StartAsync()
}
