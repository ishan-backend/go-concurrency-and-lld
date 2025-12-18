package models

import (
	"sync/atomic"
	"time"
)

var showIDGenerator int64 = 0

// Show represents a show/movie screening in a theater
type Show struct {
	ID        int64
	MovieID   int64
	TheaterID int64
	ScreenID  int64
	ShowTime  time.Time // Combined date and time
	Price     float64   // Base price for the show
}

// NewShow creates a new Show instance
func NewShow(movieID, theaterID, screenID int64, showTime time.Time, price float64) *Show {
	showID := atomic.AddInt64(&showIDGenerator, 1)
	return &Show{
		ID:        showID,
		MovieID:   movieID,
		TheaterID: theaterID,
		ScreenID:  screenID,
		ShowTime:  showTime,
		Price:     price,
	}
}

// GetShowDate returns the date part of the show time
func (s *Show) GetShowDate() time.Time {
	return time.Date(s.ShowTime.Year(), s.ShowTime.Month(), s.ShowTime.Day(), 0, 0, 0, 0, s.ShowTime.Location())
}

// GetShowTimeOnly returns the time part of the show time
func (s *Show) GetShowTimeOnly() string {
	return s.ShowTime.Format("15:04")
}

