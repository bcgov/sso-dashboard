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

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func deleteOldClientEvents() error {
	retention_period := utils.GetEnv("RETENTION_PERIOD", "1 year")

	// see https://www.postgresql.org/docs/current/datatype-datetime.html#DATATYPE-INTERVAL-INPUT
	query := "DELETE FROM client_events WHERE date < current_date - interval ?;"
	_, err := pgdb.Query(nil, query, retention_period)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func RunCronJob() {
	loc := config.LoadTimeLocation()
	cron := gocron.NewScheduler(loc)
	cron.Every(1).Day().At("02:00").Do(func() {
		deleteOldClientEvents()
	})
	cron.StartBlocking()
}
