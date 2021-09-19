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

type Comments struct {
	l *log.Logger
}

func NewComments(l *log.Logger) *Comments {
	return &Comments{l}
}

func (c *Comments) GetComments(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	c.l.Println("Handle GET Comments")

	lc, err := data.GetComments(id)
	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie doesn't exists", http.StatusBadRequest)
		return
	}

	err = lc.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (c *Comments) AddComment(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST Comment")

	comment := r.Context().Value(KeyComment{}).(data.Comment)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	commentForReturn, err := data.AddComment(&comment, reqToken)

	if err == data.ErrMovieNotFound {
		http.Error(rw, "Movie doesn't exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to create comment", http.StatusBadRequest)
		return
	}

	err = commentForReturn.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (c Comments) DeleteComment(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	c.l.Println("Handle DELETE Comment", id)

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	err = data.DeleteComment(id, reqToken)
	if err == data.ErrCommentCannotBeDeleted {
		http.Error(rw, "You can't delete someone else's comment", http.StatusBadRequest)
		return
	}

	if err == data.ErrCommentNotFound {
		http.Error(rw, "Comment not found", http.StatusBadRequest)
		return
	}
}

type KeyComment struct{}

func (c Comments) MiddlewareValidateComments(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		comment := data.Comment{}

		err := comment.FromJSON(r.Body)
		if err != nil {
			c.l.Println("[ERROR] deserializing comment", err)
			http.Error(rw, "Error reading comment", http.StatusBadRequest)
			return
		}

		err = comment.Validate()
		if err != nil {
			c.l.Println("[ERROR] validating comment", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating comment: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyComment{}, comment)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
