package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"sso-dashboard.bcgov.com/aggregator/model"
	"sso-dashboard.bcgov.com/aggregator/promtail"
	"sso-dashboard.bcgov.com/aggregator/utils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {

	root := mux.NewRouter()

	apiRouter := root.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/healthz", HealthHandler)
	apiRouter.HandleFunc("/promtail/push", promtail.PromtailPushHandler)

	port := utils.GetEnv("PORT", "8080")
	listenAddr := ":" + port

	log.Printf("attempting listen on %s", listenAddr)
	http.Handle("/", root)

	defer model.GetDB().Close()

	log.Fatalln(http.ListenAndServe(listenAddr, nil))
}
