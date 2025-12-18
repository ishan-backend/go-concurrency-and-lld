package models

import "errors"

var (
	ErrMovieNotFound        = errors.New("movie not found")
	ErrTheaterNotFound      = errors.New("theater not found")
	ErrScreenNotFound       = errors.New("screen not found")
	ErrShowNotFound         = errors.New("show not found")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrSeatNotFound         = errors.New("seat not found")
	ErrInvalidSeatPosition  = errors.New("invalid seat position")
	ErrSeatNotAvailable     = errors.New("seat not available")
	ErrSeatAlreadyBooked     = errors.New("seat already booked")
	ErrSeatLockedByOther    = errors.New("seat locked by another user")
	ErrInvalidBooking        = errors.New("invalid booking")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
)

