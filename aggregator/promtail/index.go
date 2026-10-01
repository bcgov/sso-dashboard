package promtail

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/dskit/tenant"
	"github.com/grafana/loki/v3/pkg/loghttp/push"

	util_log "github.com/grafana/loki/v3/pkg/util/log"
	"github.com/prometheus/prometheus/model/labels"
	promql_parser "github.com/prometheus/prometheus/promql/parser"

	"sso-dashboard.bcgov.com/aggregator/config"
	"sso-dashboard.bcgov.com/aggregator/model"
)

// Matches the defaults used by Loki's distributor for the push HTTP path.
const (
	maxRecvMsgSize      = 4 * 1024 * 1024
	maxDecompressedSize = 100 * 1024 * 1024
)

// Parses the idp from usernames formatted as <guid@idp>.
func idpFromUsername(username string, clientId string) string {
	// In the case of multiple or zero @ chars, there is no reliable way to parse the IDP
	if strings.Count(username, "@") != 1 {
		return "unavailable"
	}
	parts := strings.Split(username, "@")
	last := parts[len(parts)-1]
	if last == "" {
		return "unavailable"
	}
	// For the bcsc case, the clientid is the idp alias
	if last == clientId {
		return "bcsc"
	} else {
		return last
	}
}

// ParseMetric requires a Parser instance since prometheus v0.311 (was
// previously a package-level function).
var metricParser = promql_parser.NewParser(promql_parser.Options{})

func PromtailPushHandler(w http.ResponseWriter, r *http.Request) {
	logger := util_log.WithContext(r.Context(), util_log.Logger)
	userID, _ := tenant.TenantID(r.Context())
	req, _, err := push.ParseRequest(
		logger, userID, maxRecvMsgSize, maxDecompressedSize, r,
		push.EmptyLimits{}, nil, push.ParseLokiRequest, nil, nil, "", "",
	)
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
		ls, err := metricParser.ParseMetric(stream.Labels)
		if err != nil {
			lastErr = err
			continue
		}

		var (
			environment = ""
			realmId     = ""
			clientId    = ""
			eventType   = ""
			idp         = ""
			date        = time.Now()
			timestamp   = ""
			username    = ""
		)

		// labels.Labels is an opaque type in prometheus v0.311+; use Range
		// instead of directly ranging over it.
		ls.Range(func(v labels.Label) {
			switch v.Name {
			case "environment":
				environment = v.Value
			case "realmId":
				realmId = v.Value
			case "clientId":
				clientId = v.Value
			case "eventType":
				eventType = v.Value
			case "timestamp":
				timestamp = v.Value
			case "username":
				if v.Value != "" {
					username = v.Value
				}
			}
		})

		idp = idpFromUsername(username, clientId)
		t, err := time.Parse(time.RFC3339Nano, timestamp)

		if err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			continue
		}

		duration := time.Since(t)

		// only collect event logs, skip the system logs
		// reject logs older than 24 hours
		if eventType == "" || duration >= 24*time.Hour {
			continue
		}

		if environment == "" || realmId == "" || clientId == "" || eventType == "" {
			log.Printf("Skipping stream due to missing required labels")
			continue
		}

		if len(stream.Entries) == 0 {
			continue
		}

		// Alloy batches in 1 second intervals, and date is flattened to 24 hours. So just using the earliest timestamp per batch.
		entry := stream.Entries[0]
		year, month, day := entry.Timestamp.In(loc).Date()
		date = time.Date(year, month, day, 0, 0, 0, 0, loc)

		err = model.UpsertClientEventWithIDP(environment, realmId, clientId, eventType, idp, date, len(stream.Entries))
		if err != nil {
			log.Println("Upsert error:", err)
		}
	}

	if lastErr != nil {
		http.Error(w, lastErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
