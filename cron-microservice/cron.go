package main

import (
	"cron-microservice/events"
	"cron-microservice/healthz"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
)

func main() {
	gocron.Every(1).Minute().Do(events.CancelNotConfirmedReservations)

	<-gocron.Start()

	http.HandleFunc("/healthz", healthz.Index)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}