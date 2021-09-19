package data

import (
	"encoding/json"
	"io"
)

type Ticket struct {
	ID           int     `json:"id"`
	ProjectionID int     `json:"projection_id"`
	UserID 		 int     `json:"user_id"`
	MovieName    string  `json:"movie_name"`
	Customer     string  `json:"customer"`
	Row    	     int     `json:"row"`
	Column    	 int     `json:"column"`
	DateTime     float64 `json:"date_time"`
	Price        float32 `json:"price"`
	Sold         bool    `json:"sold"`
}

type Tickets []*Ticket

func (t *Ticket) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (t *Tickets) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}