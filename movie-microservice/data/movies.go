package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
)

type Movie struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Genre 		string  `json:"genre" validate:"required"`
	Length      int 	`json:"length" validate:"required,gt=0"`
	Year      	int 	`json:"year" validate:"required,gt=0,lt=2022"`
	AverageRate float32 `json:"average_rate"`
}

type Projection struct {
	ID          int     `json:"id"`
	MovieID     int 	`json:"movie_id" validate:"required"`
	MovieName   string  `json:"movie_name"`
	DateTime    float64 `json:"date_time" validate:"required"`
	Price       float32 `json:"price" validate:"required,gt=0"`
}

type Projections []*Projection

func (m *Movie) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

func (m *Movie) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func (p *Projection) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Projections) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (m *Movie) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type Movies []*Movie

func (m *Movies) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

var db *gorm.DB
var err error

func GetMovies() Movies {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var movies []Movie
	db.Find(&movies)

	var forReturn Movies
	for i, _ := range movies {
		forReturn = append(forReturn, &movies[i])
	}

	return forReturn
}

func AddMovie(m *Movie) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	m.AverageRate = 0
	db.Create(m)
}

func DeleteMovie(id int, token string) error {
	movie, err := FindMovie(id)
	if err != nil {
		return err
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9091/api/projections/by-movie/" + strconv.Itoa(id), nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var projections Projections
	if err := json.NewDecoder(resp.Body).Decode(&projections); err != nil {
		return err
	}

	if len(projections) != 0 {
		return ErrMovieCannotBeDeleted
	}

	db.Delete(movie)

	return nil
}

func UpdateMovie(id int, m *Movie, token string) error {
	movie, err := FindMovie(id)
	if err != nil {
		return err
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	oldName := movie.Name

	movie.Name = m.Name
	movie.Description = m.Description
	movie.Genre = m.Genre
	movie.Length = m.Length
	movie.Year = m.Year
	movie.AverageRate = m.AverageRate

	if oldName != m.Name {
		var bearer = "Bearer " + token
		jsonValue, _ := json.Marshal(movie)
		req, err := http.NewRequest("PUT", "http://localhost:9091/api/projections", bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	db.Save(movie)

	return nil
}

func UpdateMovieAverageRate(id int, average float32) error {
	movie, err := FindMovie(id)
	if err != nil {
		return err
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	movie.AverageRate = average

	db.Save(movie)

	return nil
}

var ErrMovieCannotBeDeleted = fmt.Errorf("Movie cannot be deleted")
var ErrMovieNotFound = fmt.Errorf("Movie not found")

func FindMovie(id int) (*Movie, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_movies sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var movie Movie
	result := db.First(&movie, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, ErrMovieNotFound
	}

	return &movie, nil
}
