package api

import (
	"net/http"

	"github.com/ypapax/status_check/config"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check"
	"github.com/ypapax/status_check/status"
)

func Serve(conf config.Config) error {
	_, statusService, err := status_check.DbServices(conf.DbType, conf.ConnectionString)
	if err != nil {
		logrus.Error(err)
		return err
	}
	statusHandler := status.NewStatusHandler(statusService)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/services-count/available/{from_ts}/{to_ts}", statusHandler.AvailableServices).Methods("GET")
	router.HandleFunc("/services-count/not-available/{from_ts}/{to_ts}", statusHandler.NotAvailableServices).Methods("GET")
	router.HandleFunc("/services-count/faster/{dur_ms}/{from_ts}/{to_ts}", statusHandler.ServicesWithResponseFasterThan).Methods("GET")
	router.HandleFunc("/services-count/slower/{dur_ms}/{from_ts}/{to_ts}", statusHandler.ServicesWithResponseSlowerThan).Methods("GET")
	http.Handle("/", accessControl(router))

	logrus.Println("Listening on  " + conf.Bind)
	if err := http.ListenAndServe(conf.Bind, nil); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
