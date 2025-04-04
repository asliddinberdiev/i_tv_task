package movie

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Year     uint16  `json:"year"`
	Genre    string  `json:"genre"`
	Rating   float32 `json:"rating"`
}
