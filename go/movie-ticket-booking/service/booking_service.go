package service

import (
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/repository"
)

const (
	SeatLockDuration = 5 * time.Minute // Lock duration for seats during booking
)

// BookingServiceImpl implements BookingService
type BookingServiceImpl struct {
	bookingRepo repository.BookingRepository
	showRepo    repository.ShowRepository
	theaterRepo repository.TheaterRepository
}

// NewBookingService creates a new BookingService instance
func NewBookingService(bookingRepo repository.BookingRepository, showRepo repository.ShowRepository, theaterRepo repository.TheaterRepository) BookingService {
	return &BookingServiceImpl{
		bookingRepo: bookingRepo,
		showRepo:    showRepo,
		theaterRepo: theaterRepo,
	}
}

// BookTickets books tickets for a show
func (s *BookingServiceImpl) BookTickets(userID, showID int64, seatIDs []string) (*models.Booking, error) {
	// Validate show exists
	show, err := s.showRepo.GetShow(showID)
	if err != nil {
		return nil, err
	}

	// Get theater and screen
	theater, err := s.theaterRepo.GetTheater(show.TheaterID)
	if err != nil {
		return nil, err
	}

	screen, err := theater.GetScreen(show.ScreenID)
	if err != nil {
		return nil, err
	}

	// Get existing bookings for this show
	existingBookings := s.bookingRepo.GetBookingsByShow(showID)
	bookedSeats := make(map[string]bool)
	for _, booking := range existingBookings {
		for _, seatID := range booking.SeatIDs {
			bookedSeats[seatID] = true
		}
	}

	// Validate and lock seats
	var seatsToBook []*models.Seat
	for _, seatID := range seatIDs {
		// Parse seat ID to get row and col
		row, col, err := parseSeatID(seatID)
		if err != nil {
			return nil, err
		}

		seat, err := screen.GetSeat(row, col)
		if err != nil {
			return nil, err
		}

		// Check if seat is already booked
		if bookedSeats[seatID] {
			return nil, models.ErrSeatAlreadyBooked
		}

		// Try to lock the seat
		if !seat.Lock(userID, SeatLockDuration) {
			return nil, models.ErrSeatNotAvailable
		}

		seatsToBook = append(seatsToBook, seat)
	}

	// Calculate total amount
	totalAmount := s.calculateTotalAmount(show, seatsToBook)

	// Create booking
	booking := models.NewBooking(userID, showID, seatIDs, totalAmount)

	// Book all seats
	for _, seat := range seatsToBook {
		if !seat.Book(userID) {
			// If booking fails, unlock all previously locked seats
			for _, s := range seatsToBook {
				s.Unlock(userID)
			}
			return nil, models.ErrSeatNotAvailable
		}
	}

	// Save booking
	err = s.bookingRepo.AddBooking(booking)
	if err != nil {
		// If save fails, unlock all seats
		for _, seat := range seatsToBook {
			seat.Unlock(userID)
		}
		return nil, err
	}

	return booking, nil
}

// CancelBooking cancels a booking
func (s *BookingServiceImpl) CancelBooking(bookingID int64) error {
	booking, err := s.bookingRepo.GetBooking(bookingID)
	if err != nil {
		return err
	}

	if booking.IsCancelled() {
		return models.ErrBookingAlreadyCancelled
	}

	// Get show and screen to unlock seats
	show, err := s.showRepo.GetShow(booking.ShowID)
	if err != nil {
		return err
	}

	theater, err := s.theaterRepo.GetTheater(show.TheaterID)
	if err != nil {
		return err
	}

	screen, err := theater.GetScreen(show.ScreenID)
	if err != nil {
		return err
	}

	// Unlock/booked seats (mark as available)
	for _, seatID := range booking.SeatIDs {
		row, col, err := parseSeatID(seatID)
		if err != nil {
			continue
		}

		seat, err := screen.GetSeat(row, col)
		if err != nil {
			continue
		}

		// Release the seat (reset status to available)
		seat.Release()
	}

	// Cancel booking
	booking.Cancel()
	return s.bookingRepo.UpdateBooking(booking)
}

// GetBooking retrieves a booking by ID
func (s *BookingServiceImpl) GetBooking(bookingID int64) (*models.Booking, error) {
	return s.bookingRepo.GetBooking(bookingID)
}

// GetUserBookings retrieves all bookings for a user
func (s *BookingServiceImpl) GetUserBookings(userID int64) []*models.Booking {
	return s.bookingRepo.GetBookingsByUser(userID)
}

// calculateTotalAmount calculates the total amount for booking
func (s *BookingServiceImpl) calculateTotalAmount(show *models.Show, seats []*models.Seat) float64 {
	basePrice := show.Price
	totalAmount := 0.0

	for _, seat := range seats {
		price := basePrice
		// Apply multipliers based on seat type
		switch seat.Type {
		case models.SeatTypeVIP:
			price *= 1.5 // VIP seats are 50% more expensive
		case models.SeatTypePremium:
			price *= 1.2 // Premium seats are 20% more expensive
		}
		totalAmount += price
	}

	return totalAmount
}

// parseSeatID parses a seat ID (e.g., "A1") to row and column indices
func parseSeatID(seatID string) (int, int, error) {
	if len(seatID) < 2 {
		return 0, 0, models.ErrInvalidSeatPosition
	}

	rowLetter := seatID[0]
	row := int(rowLetter - 'A')

	colStr := seatID[1:]
	col := 0
	for _, ch := range colStr {
		if ch < '0' || ch > '9' {
			return 0, 0, models.ErrInvalidSeatPosition
		}
		col = col*10 + int(ch-'0')
	}
	col-- // Convert to 0-based index

	return row, col, nil
}
