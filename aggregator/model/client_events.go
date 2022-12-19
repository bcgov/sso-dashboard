package model

import (
	"log"
	"time"
)

type ClientEvent struct {
	Environment string `pg:",notnull"`
	RealmID     string `pg:",notnull"`
	ClientID    string `pg:",notnull"`
	EventType   string `pg:",notnull"`
	Count       int    `pg:",notnull"`
	Date        string `pg:",notnull"`
}

func UpsertClientEvent(environment string, realmID string, clientID string, eventType string, date time.Time) (*ClientEvent, error) {
	db := ConnectDB()
	defer db.Close()

	var clientEvent ClientEvent
	query := "INSERT INTO client_events (environment, realm_id, client_id, event_type, date, count) VALUES(?, ?,?,?,?,1) ON CONFLICT (realm_id, client_id, event_type, date) DO UPDATE SET count = client_events.count + 1"
	_, err := db.Query(&clientEvent, query, environment, realmID, clientID, eventType, date)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &clientEvent, nil
}
