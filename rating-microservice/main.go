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
	"rating-microservice/data"
	"rating-microservice/handlers"
	"time"
)

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&data.Rate{})
	db.AutoMigrate(&data.Comment{})

	l := log.New(os.Stdout, "ratings-api ", log.LstdFlags)

	ch := handlers.NewComments(l)
	rh := handlers.NewRates(l)

	sm := mux.NewRouter()

	getCommentsRouter := sm.Methods(http.MethodGet).Subrouter()
	getCommentsRouter.HandleFunc("/api/comments/{id:[0-9]+}", ch.GetComments)

	postCommentRouter := sm.Methods(http.MethodPost).Subrouter()
	postCommentRouter.HandleFunc("/api/comments", ch.AddComment)
	postCommentRouter.Use(ch.MiddlewareValidateComments)

	deleteCommentRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteCommentRouter.HandleFunc("/api/comments/{id:[0-9]+}", ch.DeleteComment)

	getRateRouter := sm.Methods(http.MethodGet).Subrouter()
	getRateRouter.HandleFunc("/api/rates/{id:[0-9]+}", rh.GetRating)

	postRateRouter := sm.Methods(http.MethodPost).Subrouter()
	postRateRouter.HandleFunc("/api/rates", rh.AddRate)
	postRateRouter.Use(rh.MiddlewareValidateRates)

	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         ":9092",
		Handler:      cors(sm),
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9092")

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
