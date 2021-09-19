package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ticket-microservice/data"
)

type Tickets struct {
	l *log.Logger
}

func NewTickets(l *log.Logger) *Tickets {
	return &Tickets{l}
}

func (t *Tickets) GetPersonalTickets(rw http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle GET Personal tickets")

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	lt := data.GetPersonalTickets(reqToken)

	err := lt.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Tickets) GetReservations(rw http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle GET Reservations")

	lt := data.GetReservations()

	err := lt.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Tickets) GetSoldTickets(rw http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle GET Sold tickets")

	lt := data.GetSoldTickets()

	err := lt.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Tickets) GetByProjection(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	t.l.Println("Handle GET Tickets by projection", id)

	var tickets data.Tickets
	tickets = data.GetTicketsByProjection(id)

	errr := tickets.ToJSON(rw)
	if errr != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Tickets) AddTickets(rw http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle POST Tickets")

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	reqToken := splitToken[1]

	ticketList := r.Context().Value(KeyTicketList{}).(data.TicketList)
	err := data.AddTickets(&ticketList, reqToken)
	if err == data.ErrProjectionNotFound {
		http.Error(rw, "Projection not found", http.StatusNotFound)
		return
	}
	if err == data.ErrSeatsTaken {
		http.Error(rw, "Seats are already taken", http.StatusBadRequest)
		return
	}
}

func (t Tickets) CancelTicket(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	t.l.Println("Handle DELETE Cancel ticket", id)

	err = data.CancelTicket(id)

	if err == data.ErrTicketNotFound {
		http.Error(rw, "Ticket not found", http.StatusNotFound)
		return
	}
	if err == data.ErrCannotCancelTicket {
		http.Error(rw, "You can't cancel this ticket because it's less than 2 hours until projection", http.StatusBadRequest)
		return
	}
}

func (t Tickets) DeleteTicket(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	t.l.Println("Handle DELETE Ticket", id)

	err = data.DeleteTicket(id)

	if err == data.ErrTicketNotFound {
		http.Error(rw, "Ticket not found", http.StatusNotFound)
		return
	}
}

func (t Tickets) UpdateTicket(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	t.l.Println("Handle PUT Ticket", id)

	ticket := r.Context().Value(KeyTicket{}).(data.Ticket)

	err = data.UpdateTicket(id, &ticket)

	if err == data.ErrTicketNotFound {
		http.Error(rw, "Ticket not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Ticket not found", http.StatusInternalServerError)
		return
	}
}

type KeyTicket struct{}
type KeyTicketList struct{}

func (t Tickets) MiddlewareValidateTicket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ticket := data.Ticket{}

		err := ticket.FromJSON(r.Body)
		if err != nil {
			t.l.Println("[ERROR] deserializing ticket", err)
			http.Error(rw, "Error reading ticket", http.StatusBadRequest)
			return
		}

		err = ticket.Validate()
		if err != nil {
			t.l.Println("[ERROR] validating ticket", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating ticket: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyTicket{}, ticket)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func (t Tickets) MiddlewareValidateTicketList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ticketList := data.TicketList{}

		err := ticketList.FromJSON(r.Body)
		if err != nil {
			t.l.Println("[ERROR] deserializing ticket list", err)
			http.Error(rw, "Error reading ticket list", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyTicketList{}, ticketList)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
