package main

import (
	"log"

	"sso-dashboard.bcgov.com/aggregator/model"
)

func main() {
	log.Printf("cronjob starts...")
	model.RunEventsJob()
	model.RunSessionsJob()
	defer model.GetDB().Close()
	select {}
}
