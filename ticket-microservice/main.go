package main

import (
	"context"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ticket-microservice/data"
	"ticket-microservice/handlers"
	"time"
)

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&data.Ticket{})

	l := log.New(os.Stdout, "tickets-api ", log.LstdFlags)

	th := handlers.NewTickets(l)

	sm := mux.NewRouter()

	getPersonalTicketsRouter := sm.Methods(http.MethodGet).Subrouter()
	getPersonalTicketsRouter.HandleFunc("/api/tickets/personal", th.GetPersonalTickets)

	getReservationsRouter := sm.Methods(http.MethodGet).Subrouter()
	getReservationsRouter.HandleFunc("/api/tickets/reserved", th.GetReservations)

	getSoldTicketsRouter := sm.Methods(http.MethodGet).Subrouter()
	getSoldTicketsRouter.HandleFunc("/api/tickets/sold", th.GetSoldTickets)

	getByProjection := sm.Methods(http.MethodGet).Subrouter()
	getByProjection.HandleFunc("/api/tickets/by-projection/{id:[0-9]+}", th.GetByProjection)

	cancelRouter := sm.Methods(http.MethodDelete).Subrouter()
	cancelRouter.HandleFunc("/api/tickets/cancel/{id:[0-9]+}", th.CancelTicket)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/tickets/delete/{id:[0-9]+}", th.DeleteTicket)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/tickets", th.AddTickets)
	postRouter.Use(th.MiddlewareValidateTicketList)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/tickets/{id:[0-9]+}", th.UpdateTicket)
	putRouter.Use(th.MiddlewareValidateTicket)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         ":9093",
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9093")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
