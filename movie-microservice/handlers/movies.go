package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"movie-microservice/data"
	"net/http"
	"strconv"
	"strings"
)

type Movies struct {
	l *log.Logger
}

func NewMovies(l *log.Logger) *Movies {
	return &Movies{l}
}

func (m *Movies) GetMovies(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET Movies")

	lm := data.GetMovies()

	err := lm.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (m *Movies) GetOneMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	m.l.Println("Handle GET Movie", id)

	var movie *data.Movie
	movie, err = data.FindMovie(id)
	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Movie not found", http.StatusInternalServerError)
		return
	}

	errr := movie.ToJSON(rw)
	if errr != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (m *Movies) AddMovie(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle POST Movie")

	movie := r.Context().Value(KeyMovie{}).(data.Movie)
	data.AddMovie(&movie)
}

func (m Movies) DeleteMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	m.l.Println("Handle DELETE Movie", id)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	err = data.DeleteMovie(id, reqToken)

	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}
	if err == data.ErrMovieCannotBeDeleted {
		http.Error(rw, "Movie can't be deleted because it has projections", http.StatusBadRequest)
		return
	}
}

func (m Movies) UpdateMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	m.l.Println("Handle PUT Movie", id)

	movie := r.Context().Value(KeyMovie{}).(data.Movie)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	err = data.UpdateMovie(id, &movie, reqToken)

	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Movie not found", http.StatusInternalServerError)
		return
	}
}

func (m Movies) UpdateMovieAverageRate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	average, err := strconv.ParseFloat(vars["average"], 32)
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	m.l.Println("Handle PUT Movie", id, "average rate", average)

	err = data.UpdateMovieAverageRate(id, float32(average))

	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Movie not found", http.StatusInternalServerError)
		return
	}
}

type KeyMovie struct{}

func (m Movies) MiddlewareValidateMovie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		movie := data.Movie{}

		err := movie.FromJSON(r.Body)
		if err != nil {
			m.l.Println("[ERROR] deserializing movie", err)
			http.Error(rw, "Error reading movie", http.StatusBadRequest)
			return
		}

		err = movie.Validate()
		if err != nil {
			m.l.Println("[ERROR] validating movie", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating movie: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyMovie{}, movie)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}