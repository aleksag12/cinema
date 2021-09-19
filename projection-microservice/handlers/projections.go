package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"projection-microservice/data"
	"strconv"
	"strings"
)

type Projections struct {
	l *log.Logger
}

func NewProjections(l *log.Logger) *Projections {
	return &Projections{l}
}

func (p *Projections) GetProjections(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Projections")

	lp := data.GetProjections()

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Projections) GetOneProjection(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle GET Projection", id)

	var projection *data.Projection
	projection, err = data.FindProjection(id)
	if err == data.ErrProjectionNotFound {
		http.Error(rw, "Projection not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Projection not found", http.StatusInternalServerError)
		return
	}

	errr := projection.ToJSON(rw)
	if errr != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Projections) GetByMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle GET Projection by movie", id)

	var projections data.Projections
	projections = data.GetProjectionsByMovie(id)

	errr := projections.ToJSON(rw)
	if errr != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Projections) GetReservedSeats(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle GET Reserved seats", id)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	var seats *data.Seats
	seats, err = data.GetReservedSeats(id, reqToken)
	if err == data.ErrProjectionNotFound {
		http.Error(rw, "Projection not found", http.StatusNotFound)
	}

	errr := seats.ToJSON(rw)
	if errr != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Projections) AddProjection(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Projection")

	projection := r.Context().Value(KeyProjection{}).(data.Projection)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	err := data.AddProjection(&projection, reqToken)

	if err != nil {
		http.Error(rw, "Unable to create projection", http.StatusBadRequest)
		return
	}
}

func (p Projections) DeleteProjection(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle DELETE Projection", id)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	err = data.DeleteProjection(id, reqToken)
	if err == data.ErrProjectionNotFound {
		http.Error(rw, "Projection not found", http.StatusNotFound)
		return
	}
	if err == data.ErrProjectionCannotBeDeleted {
		http.Error(rw, "Projection can't be deleted because it has tickets", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Projection not found", http.StatusInternalServerError)
		return
	}
}

func (p *Projections) UpdateProjections(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Projections")

	movie := data.Movie{}
	movie.FromJSON(r.Body)

	data.UpdateProjections(&movie)
}

type KeyProjection struct{}

func (p Projections) MiddlewareValidateProjection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		projection := data.Projection{}

		err := projection.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing projection", err)
			http.Error(rw, "Error reading projection", http.StatusBadRequest)
			return
		}

		err = projection.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating projection", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating projection: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProjection{}, projection)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
