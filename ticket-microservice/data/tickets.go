package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Ticket struct {
	ID           int     `json:"id"`
	ProjectionID int     `json:"projection_id" validate:"required"`
	UserID 		 int     `json:"user_id"`
	MovieName    string  `json:"movie_name"`
	Customer     string  `json:"customer"`
	Row    	     int     `json:"row" validate:"required"`
	Column    	 int     `json:"column" validate:"required"`
	DateTime     float64 `json:"date_time"`
	Price        float32 `json:"price"`
	Sold         bool    `json:"sold"`
}

type TicketList struct {
	Tickets []Ticket `json:"tickets"`
}

type Projection struct {
	ID          int     `json:"id"`
	MovieID     int 	`json:"movie_id"`
	MovieName   string  `json:"movie_name"`
	DateTime    float64 `json:"date_time"`
	Price       float32 `json:"price"`
}

type Tickets []*Ticket

func (t *Ticket) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (t *Ticket) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Ticket) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

func (t *TicketList) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (p *Projection) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (t *Tickets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

var db *gorm.DB
var err error

func GetPersonalTickets(token string) Tickets {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte("UUSRHZPjPVgcIyWyGVGPp5Rj6pFaVgSg"), nil
	})

	sub := claims["sub"].(map[string]interface{})
	userID := sub["id"].(float64)
	userIDint := int(userID)

	var tickets []Ticket

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	db.Where("user_id = ? AND date_time > ?", userIDint, float64(millisec)).Find(&tickets)

	var forReturn Tickets
	for i, _ := range tickets {
		forReturn = append(forReturn, &tickets[i])
	}

	return forReturn
}

func GetReservations() Tickets {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var tickets []Ticket

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	db.Where("sold = ? AND date_time > ?", false, float64(millisec)).Find(&tickets)

	var forReturn Tickets
	for i, _ := range tickets {
		forReturn = append(forReturn, &tickets[i])
	}

	return forReturn
}

func GetSoldTickets() Tickets {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var tickets []Ticket

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	db.Where("sold = ? AND date_time > ?", true, float64(millisec)).Find(&tickets)

	var forReturn Tickets
	for i, _ := range tickets {
		forReturn = append(forReturn, &tickets[i])
	}

	return forReturn
}

func GetTicketsByProjection(id int) Tickets {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var tickets []Ticket
	db.Where(&Ticket{ProjectionID: id}).Find(&tickets)

	var forReturn Tickets
	for i, _ := range tickets {
		forReturn = append(forReturn, &tickets[i])
	}

	return forReturn
}

func AddTickets(tickets *TicketList, token string) error {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	for _, ticket := range tickets.Tickets {
		var bearer = "Bearer " + token
		req, err := http.NewRequest("GET", "http://localhost:9091/api/projections/" + strconv.Itoa(ticket.ProjectionID), nil)

		req.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return ErrProjectionNotFound
		}
		defer resp.Body.Close()

		var projection Projection
		if err := json.NewDecoder(resp.Body).Decode(&projection); err != nil {
			return ErrProjectionNotFound
		}

		ticket.MovieName = projection.MovieName
		ticket.DateTime = projection.DateTime
		ticket.Price = projection.Price

		if ticket.Customer == "" {
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
				return []byte("UUSRHZPjPVgcIyWyGVGPp5Rj6pFaVgSg"), nil
			})

			sub := claims["sub"].(map[string]interface{})
			userID := sub["id"].(float64)
			username := sub["username"].(string)

			ticket.UserID = int(userID)
			ticket.Customer = username
		}

		if takenSeat(ticket.ProjectionID, ticket.Row, ticket.Column) {
			return ErrSeatsTaken
		}

		db.Create(&ticket)
	}

	return nil
}

func takenSeat(projectionID int, row int, column int) bool {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")

	var t Ticket
	result := db.Where(&Ticket{ProjectionID: projectionID, Row: row, Column: column}).Find(&t)

	if result.RowsAffected == 0 {
		return false
	}

	return true
}

func CancelTicket(id int) error {
	ticket, err := FindTicket(id)
	if err != nil {
		return ErrTicketNotFound
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	if float64(millisec) > (ticket.DateTime - float64(7200000)) {
		return ErrCannotCancelTicket
	}

	db.Delete(ticket)

	return nil
}

func DeleteTicket(id int) error {
	ticket, err := FindTicket(id)
	if err != nil {
		return ErrTicketNotFound
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Delete(ticket)

	return nil
}

func UpdateTicket(id int, t *Ticket) error {
	ticket, err := FindTicket(id)
	if err != nil {
		return ErrTicketNotFound
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	ticket.Sold = t.Sold

	db.Save(ticket)

	return nil
}

var ErrProjectionNotFound = fmt.Errorf("Projection not found")
var ErrTicketNotFound = fmt.Errorf("Ticket not found")
var ErrSeatsTaken = fmt.Errorf("Seats are already taken")
var ErrCannotCancelTicket = fmt.Errorf("You can't cancel this ticket because it's less than 2 hours until projection")

func FindTicket(id int) (*Ticket, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_tickets sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var ticket Ticket
	result := db.First(&ticket, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, ErrTicketNotFound
	}

	return &ticket, nil
}
