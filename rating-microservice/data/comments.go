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
)

type Comment struct {
	ID          int     `json:"id"`
	MovieID     int 	`json:"movie_id" validate:"required"`
	UserID     	int 	`json:"user_id""`
	Username    string  `json:"username"`
	Text    	string  `json:"text" validate:"required"`
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

func (m *Movie) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

func (m *Movie) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

type Comments []*Comment

func (c *Comment) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}

func (c *Comment) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *Comments) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

var db *gorm.DB
var err error

func GetComments(id int) (Comments, error) {
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

	var comments []Comment
	db.Where(&Comment{MovieID: id}).Find(&comments)

	var forReturn Comments
	for i, _ := range comments {
		forReturn = append(forReturn, &comments[i])
	}

	return forReturn, nil
}

func AddComment(c *Comment, token string) (*Comment, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://localhost:9090/api/movies/" + strconv.Itoa(c.MovieID), nil)

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
	username := sub["username"].(string)

	c.UserID = int(id)
	c.Username = username

	db.Create(c)

	return c, nil
}

func DeleteComment(id int, token string) error {
	comment, err := findComment(id)
	if err == ErrCommentNotFound {
		return err
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")
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

	if userIDint != comment.UserID {
		return ErrCommentCannotBeDeleted
	}

	db.Delete(comment)

	return nil
}

var ErrCommentCannotBeDeleted = fmt.Errorf("You can't delete someone else's comment")
var ErrCommentNotFound = fmt.Errorf("Comment not found")
var ErrMovieNotFound = fmt.Errorf("Movie not found")

func findComment(id int) (*Comment, error) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go_cinema_ratings sslmode=disable password=dejanradonjic")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var comment Comment
	result := db.First(&comment, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, ErrCommentNotFound
	}

	return &comment, nil
}