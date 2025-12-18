package models

import (
	"sync/atomic"
	"time"
)

var bookingIDGenerator int64 = 0

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

// Booking represents a ticket booking
type Booking struct {
	ID          int64
	UserID      int64
	ShowID      int64
	SeatIDs     []string // List of seat IDs booked
	TotalAmount float64
	Status      BookingStatus
	BookedAt    time.Time
}

// NewBooking creates a new Booking instance
func NewBooking(userID, showID int64, seatIDs []string, totalAmount float64) *Booking {
	bookingID := atomic.AddInt64(&bookingIDGenerator, 1)
	return &Booking{
		ID:          bookingID,
		UserID:      userID,
		ShowID:      showID,
		SeatIDs:     seatIDs,
		TotalAmount: totalAmount,
		Status:      BookingStatusConfirmed,
		BookedAt:    time.Now(),
	}
}

// Cancel cancels the booking
func (b *Booking) Cancel() {
	b.Status = BookingStatusCancelled
}

// IsCancelled checks if the booking is cancelled
func (b *Booking) IsCancelled() bool {
	return b.Status == BookingStatusCancelled
}

