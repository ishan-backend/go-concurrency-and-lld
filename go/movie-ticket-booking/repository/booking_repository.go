package repository

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
)

// BookingRepositoryImpl implements BookingRepository
type BookingRepositoryImpl struct {
	bookings map[int64]*models.Booking
	mu       sync.RWMutex
}

// NewBookingRepository creates a new BookingRepository instance
func NewBookingRepository() BookingRepository {
	return &BookingRepositoryImpl{
		bookings: make(map[int64]*models.Booking),
	}
}

// AddBooking adds a booking to the repository
func (r *BookingRepositoryImpl) AddBooking(booking *models.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[booking.ID] = booking
	return nil
}

// GetBooking retrieves a booking by ID
func (r *BookingRepositoryImpl) GetBooking(bookingID int64) (*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	booking, exists := r.bookings[bookingID]
	if !exists {
		return nil, models.ErrBookingNotFound
	}
	return booking, nil
}

// GetBookingsByUser retrieves all bookings for a user
func (r *BookingRepositoryImpl) GetBookingsByUser(userID int64) []*models.Booking {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Booking
	for _, booking := range r.bookings {
		if booking.UserID == userID {
			results = append(results, booking)
		}
	}

	return results
}

// GetBookingsByShow retrieves all bookings for a show
func (r *BookingRepositoryImpl) GetBookingsByShow(showID int64) []*models.Booking {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Booking
	for _, booking := range r.bookings {
		if booking.ShowID == showID && !booking.IsCancelled() {
			results = append(results, booking)
		}
	}

	return results
}

// UpdateBooking updates a booking in the repository
func (r *BookingRepositoryImpl) UpdateBooking(booking *models.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[booking.ID] = booking
	return nil
}
