package promtail

import (
	"log"
	"net/http"
	"time"

	"github.com/grafana/dskit/tenant"
	"github.com/grafana/loki/pkg/loghttp/push"

	util_log "github.com/grafana/loki/pkg/util/log"
	promql_parser "github.com/prometheus/prometheus/promql/parser"

	"sso-dashboard.bcgov.com/aggregator/model"
)

func PromtailPushHandler(w http.ResponseWriter, r *http.Request) {
	logger := util_log.WithContext(r.Context(), util_log.Logger)
	userID, _ := tenant.TenantID(r.Context())
	req, err := push.ParseRequest(logger, userID, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var lastErr error
	for _, stream := range req.Streams {
		ls, err := promql_parser.ParseMetric(stream.Labels)
		if err != nil {
			lastErr = err
			continue
		}

		var (
			environment = ""
			realmId     = ""
			clientId    = ""
			eventType   = ""
			date        = time.Now()
		)

		for _, v := range ls {
			switch v.Name {
			case "environment":
				environment = v.Value
			case "realm_id":
				realmId = v.Value
			case "client_id":
				clientId = v.Value
			case "event_type":
				eventType = v.Value
			}
		}

		for _, entry := range stream.Entries {
			year, month, day := entry.Timestamp.Date()
			date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}

		log.Println(environment, realmId, clientId, eventType, date)
		model.UpsertClientEvent(environment, realmId, clientId, eventType, date)
	}

	if lastErr != nil {
		http.Error(w, lastErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
