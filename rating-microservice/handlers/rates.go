package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rating-microservice/data"
	"strconv"
	"strings"
)

type Rates struct {
	l *log.Logger
}

func NewRates(l *log.Logger) *Rates {
	return &Rates{l}
}

func (rat *Rates) GetRating(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	rat.l.Println("Handle GET Rate for movie", movieID)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	rate, err := data.GetRate(movieID, reqToken)
	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie doesn't exists", http.StatusBadRequest)
		return
	}

	err = rate.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (rat *Rates) AddRate(rw http.ResponseWriter, r *http.Request) {
	rat.l.Println("Handle POST Rate")

	rate := r.Context().Value(KeyRate{}).(data.Rate)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	movie, err := data.AddRate(&rate, reqToken)

	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie doesn't exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to create comment", http.StatusBadRequest)
		return
	}

	err = movie.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

type KeyRate struct{}

func (rat Rates) MiddlewareValidateRates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rate := data.Rate{}

		err := rate.FromJSON(r.Body)
		if err != nil {
			rat.l.Println("[ERROR] deserializing rate", err)
			http.Error(rw, "Error reading rate", http.StatusBadRequest)
			return
		}

		err = rate.Validate()
		if err != nil {
			rat.l.Println("[ERROR] validating rate", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating rate: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyRate{}, rate)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
