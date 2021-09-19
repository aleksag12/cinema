package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Projection struct {
	ID          int     `json:"id"`
	MovieID     int 	`json:"movie_id" validate:"required"`
	MovieName   string  `json:"movie_name"`
	DateTime    float64 `json:"date_time" validate:"required"`
	Price       float32 `json:"price" validate:"required,gt=0"`
}

type Movie struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Genre 		string  `json:"genre"`
	Length      int 	`json:"length"`
	Year      	int 	`json:"year"`
	AverageRate float32 `json:"average_rate"`
}

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

type Seats struct {
	Seats []string `json:"seats"`
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

func (m *Movie) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

func (p *Projection) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Projection) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (s *Seats) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

func (p *Projection) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type Projections []*Projection

func (p *Projections) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var db *gorm.DB
var err error

func GetProjections() Projections {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var projections []Projection

	now := time.Now()
	unixNano := now.UnixNano()
	millisec := unixNano / 1000000

	db.Where("date_time > ?", float64(millisec)).Find(&projections)

	var forReturn Projections
	for i, _ := range projections {
		forReturn = append(forReturn, &projections[i])
	}

	return forReturn
}

func GetProjectionsByMovie(id int) Projections {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var projections []Projection
	db.Where(&Projection{MovieID: id}).Find(&projections)

	var forReturn Projections
	for i, _ := range projections {
		forReturn = append(forReturn, &projections[i])
	}

	return forReturn
}

func GetReservedSeats(id int, token string) (*Seats, error) {
	_, err := FindProjection(id)
	if err != nil {
		return nil, ErrProjectionNotFound
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9093/api/tickets/by-projection/" + strconv.Itoa(id), nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickets Tickets
	if err := json.NewDecoder(resp.Body).Decode(&tickets); err != nil {
		return nil, err
	}

	var seats Seats
	seats.Seats = []string{}

	var stringSeat string
	for _, p := range tickets {
		stringSeat = strconv.Itoa(p.Row) + "." + strconv.Itoa(p.Column)
		seats.Seats = append(seats.Seats, stringSeat)
	}

	return &seats, nil
}

func AddProjection(p *Projection, token string) error {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9090/api/movies/" + strconv.Itoa(p.MovieID), nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return err
	}

	p.MovieName = movie.Name
	db.Create(p)

	return nil
}

func UpdateProjections(m *Movie) error {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var projs []Projection
	db.Where(&Projection{MovieID: m.ID}).Find(&projs)

	for _, p := range projs {
		p.MovieName = m.Name
		db.Save(p)
	}

	return nil
}

func DeleteProjection(id int, token string) error {
	projection, err := FindProjection(id)
	if err != nil {
		return ErrProjectionNotFound
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9093/api/tickets/by-projection/" + strconv.Itoa(id), nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tickets Tickets
	if err := json.NewDecoder(resp.Body).Decode(&tickets); err != nil {
		return err
	}

	if len(tickets) != 0 {
		return ErrProjectionCannotBeDeleted
	}

	db.Delete(projection)

	return nil
}

var ErrProjectionNotFound = fmt.Errorf("Projection not found")
var ErrProjectionCannotBeDeleted = fmt.Errorf("Projection can't be deleted because it has tickets")

func FindProjection(id int) (*Projection, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_projections sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var projection Projection
	result := db.First(&projection, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, ErrProjectionNotFound
	}

	return &projection, nil
}
