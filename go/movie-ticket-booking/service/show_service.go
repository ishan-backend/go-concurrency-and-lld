package service

import (
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/repository"
)

// ShowServiceImpl implements ShowService
type ShowServiceImpl struct {
	showRepo    repository.ShowRepository
	theaterRepo repository.TheaterRepository
	movieRepo   repository.MovieRepository
	bookingRepo repository.BookingRepository
}

// NewShowService creates a new ShowService instance
func NewShowService(showRepo repository.ShowRepository, theaterRepo repository.TheaterRepository, movieRepo repository.MovieRepository, bookingRepo repository.BookingRepository) ShowService {
	return &ShowServiceImpl{
		showRepo:    showRepo,
		theaterRepo: theaterRepo,
		movieRepo:   movieRepo,
		bookingRepo: bookingRepo,
	}
}

// CreateShow creates a new show
func (s *ShowServiceImpl) CreateShow(movieID, theaterID, screenID int64, showTime string, price float64) (*models.Show, error) {
	// Validate movie exists
	_, err := s.movieRepo.GetMovie(movieID)
	if err != nil {
		return nil, err
	}

	// Validate theater exists
	theater, err := s.theaterRepo.GetTheater(theaterID)
	if err != nil {
		return nil, err
	}

	// Validate screen exists in theater
	_, err = theater.GetScreen(screenID)
	if err != nil {
		return nil, err
	}

	// Parse show time (format: "2006-01-02 15:04")
	showDateTime, err := time.Parse("2006-01-02 15:04", showTime)
	if err != nil {
		return nil, err
	}

	show := models.NewShow(movieID, theaterID, screenID, showDateTime, price)
	err = s.showRepo.AddShow(show)
	if err != nil {
		return nil, err
	}

	return show, nil
}

// GetShow retrieves a show by ID
func (s *ShowServiceImpl) GetShow(showID int64) (*models.Show, error) {
	return s.showRepo.GetShow(showID)
}

// GetShowsByMovieAndCity retrieves shows for a movie in a specific city
func (s *ShowServiceImpl) GetShowsByMovieAndCity(movieID int64, city string) []*models.Show {
	// Get all theaters in the city
	theaters := s.theaterRepo.GetTheatersByCity(city)
	theaterIDs := make(map[int64]bool)
	for _, theater := range theaters {
		theaterIDs[theater.ID] = true
	}

	// Get all shows for the movie
	allShows := s.showRepo.GetShowsByMovieAndCity(movieID, city)

	// Filter shows by theater city
	var results []*models.Show
	for _, show := range allShows {
		if theaterIDs[show.TheaterID] {
			results = append(results, show)
		}
	}

	return results
}

// GetAvailableSeats retrieves the seat layout for a show
func (s *ShowServiceImpl) GetAvailableSeats(showID int64) ([][]*models.Seat, error) {
	show, err := s.showRepo.GetShow(showID)
	if err != nil {
		return nil, err
	}

	theater, err := s.theaterRepo.GetTheater(show.TheaterID)
	if err != nil {
		return nil, err
	}

	screen, err := theater.GetScreen(show.ScreenID)
	if err != nil {
		return nil, err
	}

	// Get all bookings for this show to mark seats as booked
	bookings := s.bookingRepo.GetBookingsByShow(showID)

	bookedSeats := make(map[string]bool)
	for _, booking := range bookings {
		for _, seatID := range booking.SeatIDs {
			bookedSeats[seatID] = true
		}
	}

	// Create a copy of the seat layout with updated status
	seatLayout := make([][]*models.Seat, screen.TotalRows)
	for i := 0; i < screen.TotalRows; i++ {
		seatLayout[i] = make([]*models.Seat, screen.TotalCols)
		for j := 0; j < screen.TotalCols; j++ {
			seat := screen.SeatLayout[i][j]
			// Create a copy to avoid modifying the original
			seatCopy := &models.Seat{
				ID:       seat.ID,
				ScreenID: seat.ScreenID,
				Row:      seat.Row,
				Col:      seat.Col,
				Type:     seat.Type,
				Status:   seat.Status,
			}

			// Update status if seat is booked
			if bookedSeats[seat.ID] {
				seatCopy.Status = models.SeatStatusBooked
			}

			seatLayout[i][j] = seatCopy
		}
	}

	return seatLayout, nil
}
