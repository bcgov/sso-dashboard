package model

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"sso-dashboard.bcgov.com/aggregator/config"
	"sso-dashboard.bcgov.com/aggregator/utils"
)

type ClientEvent struct {
	Environment string `pg:",notnull"`
	RealmID     string `pg:",notnull"`
	ClientID    string `pg:",notnull"`
	EventType   string `pg:",notnull"`
	Count       int    `pg:",notnull"`
	Date        string `pg:",notnull"`
}

func UpsertClientEvent(environment string, realmID string, clientID string, eventType string, date time.Time) error {
	query := "INSERT INTO client_events (environment, realm_id, client_id, event_type, date, count) VALUES(?,?,?,?,?,1) ON CONFLICT (environment, realm_id, client_id, event_type, date) DO UPDATE SET count = client_events.count + 1"
	_, err := pgdb.Query(nil, query, environment, realmID, clientID, eventType, date)
	defer pgdb.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func deleteOldClientEvents() error {
	retention_period := utils.GetEnv("RETENTION_PERIOD", "1 year")

	// see https://www.postgresql.org/docs/current/datatype-datetime.html#DATATYPE-INTERVAL-INPUT
	eventsQuery := "DELETE FROM client_events WHERE date < current_date - interval ?;"

	// Using same retention period for sessions stats. May need to change that in the future.
	sessionsQuery := "DELETE FROM client_sessions WHERE date < current_timestamp - interval ?;"

	_, eventsErr := pgdb.Query(nil, eventsQuery, retention_period)
	_, sessionsErr := pgdb.Query(nil, sessionsQuery, retention_period)

	if eventsErr != nil {
		log.Println(eventsErr)
		return eventsErr
	}

	if sessionsErr != nil {
		log.Println(sessionsErr)
		return sessionsErr
	}
	return nil
}

func RunEventsJob() {
	loc := config.LoadTimeLocation()
	cron := gocron.NewScheduler(loc)
	cron.Every(1).Day().At("02:00").Do(func() {
		deleteOldClientEvents()
	})
	cron.StartAsync()
}
