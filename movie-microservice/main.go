package main

import (
	"context"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"movie-microservice/data"
	"movie-microservice/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var db *gorm.DB
var err error

var (
	movies = []data.Movie{
		{
			Name:        "Movie 1",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
			Genre:       "Action",
			Length:      120,
			Year:   	 2021,
			AverageRate: 0,
		},
		{
			Name:        "Movie 2",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
			Genre:       "Comedy",
			Length:      90,
			Year:   	 2015,
			AverageRate: 0,
		},
	}
)

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&data.Movie{})

	for index := range movies {
		var existingName []data.Movie
		result := db.Where(&data.Movie{Name: movies[index].Name}).Find(&existingName)
		if result.RowsAffected == 0 {
			db.Create(&movies[index])
		}
	}

	l := log.New(os.Stdout, "movies-api ", log.LstdFlags)

	mh := handlers.NewMovies(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/movies", mh.GetMovies)

	getOneRouter := sm.Methods(http.MethodGet).Subrouter()
	getOneRouter.HandleFunc("/api/movies/{id:[0-9]+}", mh.GetOneMovie)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/movies/{id:[0-9]+}", mh.DeleteMovie)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/movies", mh.AddMovie)
	postRouter.Use(mh.MiddlewareValidateMovie)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/movies/{id:[0-9]+}", mh.UpdateMovie)
	putRouter.Use(mh.MiddlewareValidateMovie)

	updateAverageRateRouter := sm.Methods(http.MethodGet).Subrouter()
	updateAverageRateRouter.HandleFunc("/api/movies/update-average-rate/{id:[0-9]+}/{average[0-9]+\\.?[0-9]*}", mh.UpdateMovieAverageRate)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

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
