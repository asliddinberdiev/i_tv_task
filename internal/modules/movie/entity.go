package movie

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title     string         `gorm:"type:varchar(255);not null;unique;index"`
	Year      int            `gorm:"type:int;not null;index"`
	Genre     string         `gorm:"type:varchar(255);not null;index"`
	Rating    float64        `gorm:"type:float;not null;index"`
	Director  string         `gorm:"type:varchar(255);not null;index"`
}

type MovieResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Genre     string    `json:"genre"`
	Rating    float64   `json:"rating"`
	Director  string    `json:"director"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieCreateInput struct {
	Title    string  `json:"title" validate:"required,min=2,lowercase"`
	Year     int     `json:"year" validate:"required,min=1800"`
	Genre    string  `json:"genre" validate:"required,min=2,lowercase"`
	Rating   float64 `json:"rating" validate:"required,min=0,max=10"`
	Director string  `json:"director" validate:"required,min=2,lowercase"`
}

type MovieUpdateInput struct {
	ID       uint    `json:"-"`
	Title    string  `json:"title" validate:"required,min=2,lowercase"`
	Year     int     `json:"year" validate:"required,min=1800"`
	Genre    string  `json:"genre" validate:"required,min=2,lowercase"`
	Rating   float64 `json:"rating" validate:"required,min=0,max=10"`
	Director string  `json:"director" validate:"required,min=2,lowercase"`
}

type MovieListResponse struct {
	Movies []MovieResponse `json:"movies"`
	Total  uint64          `json:"total"`
}
