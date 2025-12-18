package models

import (
	"sync/atomic"
	"time"
)

var movieIDGenerator int64 = 0

// Movie represents a movie in the system
type Movie struct {
	ID          int64
	Title       string
	Duration    int // Duration in minutes
	Genre       string
	Language    string
	ReleaseDate time.Time
}

// NewMovie creates a new Movie instance
func NewMovie(title string, duration int, genre, language string, releaseDate time.Time) *Movie {
	movieID := atomic.AddInt64(&movieIDGenerator, 1)
	return &Movie{
		ID:          movieID,
		Title:       title,
		Duration:    duration,
		Genre:       genre,
		Language:    language,
		ReleaseDate: releaseDate,
	}
}
