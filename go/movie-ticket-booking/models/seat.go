package models

import (
	"sync"
	"time"
)

// SeatStatus represents the status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusBooked    SeatStatus = "BOOKED"
	SeatStatusLocked    SeatStatus = "LOCKED"
)

// Seat represents a seat in a screen
type Seat struct {
	ID       string // Format: "RowCol" e.g., "A1", "B5"
	ScreenID int64
	Row      int
	Col      int
	Type     SeatType
	Status   SeatStatus
	LockedBy int64     // User ID who locked the seat
	LockedAt time.Time // When the seat was locked
	mu       sync.RWMutex
}

// NewSeat creates a new Seat instance
func NewSeat(screenID int64, row, col int, seatType SeatType) *Seat {
	seatID := generateSeatID(row, col)
	return &Seat{
		ID:       seatID,
		ScreenID: screenID,
		Row:      row,
		Col:      col,
		Type:     seatType,
		Status:   SeatStatusAvailable,
	}
}

// generateSeatID generates a seat ID from row and column
func generateSeatID(row, col int) string {
	// Convert row number to letter (0 -> A, 1 -> B, etc.)
	rowLetter := string(rune('A' + row))
	return rowLetter + string(rune('0'+(col+1)))
}

// Lock locks the seat for a user
func (s *Seat) Lock(userID int64, lockDuration time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status != SeatStatusAvailable {
		return false
	}

	s.Status = SeatStatusLocked
	s.LockedBy = userID
	s.LockedAt = time.Now()

	// Auto-unlock after lock duration
	go func() {
		time.Sleep(lockDuration)
		s.mu.Lock()
		if s.Status == SeatStatusLocked && s.LockedBy == userID {
			s.Status = SeatStatusAvailable
			s.LockedBy = 0
		}
		s.mu.Unlock()
	}()

	return true
}

// Book books the seat
func (s *Seat) Book(userID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status == SeatStatusBooked {
		return false
	}

	// Can book if available or locked by the same user
	if s.Status == SeatStatusAvailable || (s.Status == SeatStatusLocked && s.LockedBy == userID) {
		s.Status = SeatStatusBooked
		s.LockedBy = 0
		return true
	}

	return false
}

// Unlock unlocks the seat
func (s *Seat) Unlock(userID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status == SeatStatusLocked && s.LockedBy == userID {
		s.Status = SeatStatusAvailable
		s.LockedBy = 0
		return true
	}

	return false
}

// IsAvailable checks if the seat is available
func (s *Seat) IsAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status == SeatStatusAvailable
}

// IsLockedBy checks if the seat is locked by a specific user
func (s *Seat) IsLockedBy(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status == SeatStatusLocked && s.LockedBy == userID
}

// Release releases a booked seat and makes it available again
func (s *Seat) Release() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status == SeatStatusBooked {
		s.Status = SeatStatusAvailable
		s.LockedBy = 0
		return true
	}

	return false
}
