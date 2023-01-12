package config

import (
	"sso-dashboard.bcgov.com/aggregator/utils"
	"time"
)

var (
	loc *time.Location
)

func init() {
	tz := utils.GetEnv("TZ", "")

	if tz != "" {
		loc, _ = time.LoadLocation(tz)
	} else {
		loc = time.Local
	}
}

func LoadTimeLocation() *time.Location {
	return loc
}
