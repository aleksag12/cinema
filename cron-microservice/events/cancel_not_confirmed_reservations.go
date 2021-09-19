package events

import (
	"cron-microservice/data"
	"encoding/json"
	"errors"
	"github.com/codeship/go-retro"
	"log"
	"net/http"
	"strconv"
	"time"
)

var ErrNetwork = retro.NewStaticRetryableError(errors.New("error: failed to connect"), 5, 10)

func CancelNotConfirmedReservations() {
	err := retro.DoWithRetry(func() error {
		return cancel()
	})

	if err != nil {
		log.Fatal("Failed to send request %s\n", err.Error())
		return
	}

	log.Print("Successfully sent request")
}

func cancel() error {
	req, err := http.NewRequest("GET", "http://localhost:9093/api/tickets/reserved", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tickets data.Tickets
	if err := json.NewDecoder(resp.Body).Decode(&tickets); err != nil {
		return err
	}

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	for _, t := range tickets {
		if (t.DateTime - float64(millisec)) < 1800000 {
			req, err := http.NewRequest("DELETE", "http://localhost:9093/api/tickets/delete/" + strconv.Itoa(t.ID), nil)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
		}
	}

	return nil
}