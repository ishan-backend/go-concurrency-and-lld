package models

import (
	"sync/atomic"
)

var screenIDGenerator int64 = 0

// SeatType represents the type of seat
type SeatType string

const (
	SeatTypeNormal  SeatType = "NORMAL"
	SeatTypePremium SeatType = "PREMIUM"
	SeatTypeVIP     SeatType = "VIP"
)

// Screen represents a screen in a theater
type Screen struct {
	ID         int64
	TheaterID  int64
	ScreenNum  int
	TotalRows  int
	TotalCols  int
	SeatLayout [][]*Seat // 2D array representing seat layout
}

// NewScreen creates a new Screen instance
func NewScreen(theaterID int64, screenNum, totalRows, totalCols int) *Screen {
	screenID := atomic.AddInt64(&screenIDGenerator, 1)
	
	// Initialize seat layout
	seatLayout := make([][]*Seat, totalRows)
	for i := 0; i < totalRows; i++ {
		seatLayout[i] = make([]*Seat, totalCols)
		for j := 0; j < totalCols; j++ {
			// Determine seat type based on row (first 2 rows VIP, next 3 Premium, rest Normal)
			var seatType SeatType = SeatTypeNormal
			if i < 2 {
				seatType = SeatTypeVIP
			} else if i < 5 {
				seatType = SeatTypePremium
			}
			
			seatLayout[i][j] = NewSeat(screenID, i, j, seatType)
		}
	}
	
	return &Screen{
		ID:         screenID,
		TheaterID:  theaterID,
		ScreenNum:  screenNum,
		TotalRows:  totalRows,
		TotalCols:  totalCols,
		SeatLayout: seatLayout,
	}
}

// GetSeat returns a seat at the given row and column
func (s *Screen) GetSeat(row, col int) (*Seat, error) {
	if row < 0 || row >= s.TotalRows || col < 0 || col >= s.TotalCols {
		return nil, ErrInvalidSeatPosition
	}
	return s.SeatLayout[row][col], nil
}

