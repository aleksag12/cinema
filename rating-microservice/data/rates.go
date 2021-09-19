package data

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
)

type Rate struct {
	ID          int     `json:"id"`
	MovieID     int 	`json:"movie_id" validate:"required"`
	UserID     	int 	`json:"user_id""`
	Value    	int  	`json:"value" validate:"required,gt=0,lt=6"`
}

func (rate *Rate) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(rate)
}

func (rate *Rate) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rate)
}

func (rate *Rate) Validate() error {
	validate := validator.New()
	return validate.Struct(rate)
}

func GetRate(id int, token string) (*Rate, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	req, err := http.NewRequest("GET", "http://localhost:9090/api/movies/" + strconv.Itoa(id), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrMovieNotFound
	}
	defer resp.Body.Close()

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, ErrMovieNotFound
	}

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte("UUSRHZPjPVgcIyWyGVGPp5Rj6pFaVgSg"), nil
	})

	sub := claims["sub"].(map[string]interface{})
	userID := sub["id"].(float64)
	userIDint := int(userID)

	var rate Rate
	db.Where(&Rate{MovieID: id, UserID: userIDint}).Find(&rate)

	return &rate, nil
}

func AddRate(r *Rate, token string) (*Movie, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9090/api/movies/" + strconv.Itoa(r.MovieID), nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrMovieNotFound
	}
	defer resp.Body.Close()

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, ErrMovieNotFound
	}

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte("UUSRHZPjPVgcIyWyGVGPp5Rj6pFaVgSg"), nil
	})

	sub := claims["sub"].(map[string]interface{})
	id := sub["id"].(float64)

	r.UserID = int(id)

	var rate Rate
	result := db.Where(&Rate{MovieID: r.MovieID, UserID: int(id)}).Find(&rate)

	if result.RowsAffected == 0 {
		db.Create(r)
	} else {
		rate.Value = r.Value
		db.Save(rate)
	}

	var rates []Rate
	db.Where(&Rate{MovieID: r.MovieID}).Find(&rates)

	var sum float32 = 0
	for _, r := range rates {
		sum = sum + float32(r.Value)
	}

	movie.AverageRate = sum/float32(len(rates))

	jsonValue, _ := json.Marshal(movie)
	putReq, err := http.NewRequest("PUT", "http://localhost:9090/api/movies/" + strconv.Itoa(r.MovieID), bytes.NewBuffer(jsonValue))

	putReq.Header.Add("Authorization", bearer)

	client = &http.Client{}
	resp, err = client.Do(putReq)
	if err != nil {
		return nil, ErrMovieNotFound
	}
	defer resp.Body.Close()

	return &movie, nil
}
