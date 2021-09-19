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
	"projection-microservice/data"
	"projection-microservice/handlers"
	"time"
)

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&data.Projection{})

	l := log.New(os.Stdout, "projections-api ", log.LstdFlags)

	ph := handlers.NewProjections(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/projections", ph.GetProjections)

	getOneRouter := sm.Methods(http.MethodGet).Subrouter()
	getOneRouter.HandleFunc("/api/projections/{id:[0-9]+}", ph.GetOneProjection)

	getByMovieRouter := sm.Methods(http.MethodGet).Subrouter()
	getByMovieRouter.HandleFunc("/api/projections/by-movie/{id:[0-9]+}", ph.GetByMovie)

	getSeatsRouter := sm.Methods(http.MethodGet).Subrouter()
	getSeatsRouter.HandleFunc("/api/seats/{id:[0-9]+}", ph.GetReservedSeats)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/projections/{id:[0-9]+}", ph.DeleteProjection)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/projections", ph.AddProjection)
	postRouter.Use(ph.MiddlewareValidateProjection)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/projections", ph.UpdateProjections)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         ":9091",
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9091")

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
