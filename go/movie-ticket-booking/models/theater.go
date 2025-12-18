package models

import (
	"sync/atomic"
)

var theaterIDGenerator int64 = 0

// Location represents the location of a theater
type Location struct {
	City    string
	Address string
}

// Theater represents a theater in the system
type Theater struct {
	ID       int64
	Name     string
	Location Location
	Screens  map[int64]*Screen // Map of screenID -> Screen
}

// NewTheater creates a new Theater instance
func NewTheater(name, city, address string) *Theater {
	theaterID := atomic.AddInt64(&theaterIDGenerator, 1)
	return &Theater{
		ID:   theaterID,
		Name: name,
		Location: Location{
			City:    city,
			Address: address,
		},
		Screens: make(map[int64]*Screen),
	}
}

// AddScreen adds a screen to the theater
func (t *Theater) AddScreen(screen *Screen) {
	t.Screens[screen.ID] = screen
}

// GetScreen returns a screen by ID
func (t *Theater) GetScreen(screenID int64) (*Screen, error) {
	screen, exists := t.Screens[screenID]
	if !exists {
		return nil, ErrScreenNotFound
	}
	return screen, nil
}
