package promtail

import (
	"net/http"
	"time"

	"github.com/grafana/dskit/tenant"
	"github.com/grafana/loki/pkg/loghttp/push"

	util_log "github.com/grafana/loki/pkg/util/log"
	promql_parser "github.com/prometheus/prometheus/promql/parser"

	"sso-dashboard.bcgov.com/aggregator/config"
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

	loc := config.LoadTimeLocation()

	var lastErr error

	/*
		Each unique label combination will have a dedicated stream. The logs associated to that label combination will be in stream.Entries.
		e.g. if 4 LOGIN events for the same client/realm/env are sent, there will be 1 stream with 4 entries.
		If multiple label combinations are sent in the same batch they will be in different streams.
	*/
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

		// only collect event logs, skip the system logs
		if eventType == "" {
			continue
		}

		// For the aggregated count, timestamp is flattened to the date and the count of entries in this batch will be added to the total.
		for _, entry := range stream.Entries {
			year, month, day := entry.Timestamp.In(loc).Date()
			date = time.Date(year, month, day, 0, 0, 0, 0, loc)
		}

		go model.UpsertClientEvent(environment, realmId, clientId, eventType, date, len(stream.Entries))
	}

	if lastErr != nil {
		http.Error(w, lastErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
