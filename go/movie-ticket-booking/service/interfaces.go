package service

import "github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"

// MovieService defines the interface for movie operations
type MovieService interface {
	AddMovie(title string, duration int, genre, language string, releaseDate string) (*models.Movie, error)
	GetMovie(movieID int64) (*models.Movie, error)
	SearchMoviesByTitle(title string) []*models.Movie
	SearchMoviesByGenre(genre string) []*models.Movie
	SearchMoviesByLanguage(language string) []*models.Movie
}

// TheaterService defines the interface for theater operations
type TheaterService interface {
	AddTheater(name, city, address string) (*models.Theater, error)
	GetTheater(theaterID int64) (*models.Theater, error)
	AddScreen(theaterID int64, screenNum, totalRows, totalCols int) (*models.Screen, error)
	GetTheatersByCity(city string) []*models.Theater
}

// ShowService defines the interface for show operations
type ShowService interface {
	CreateShow(movieID, theaterID, screenID int64, showTime string, price float64) (*models.Show, error)
	GetShow(showID int64) (*models.Show, error)
	GetShowsByMovieAndCity(movieID int64, city string) []*models.Show
	GetAvailableSeats(showID int64) ([][]*models.Seat, error)
}

// BookingService defines the interface for booking operations
type BookingService interface {
	BookTickets(userID, showID int64, seatIDs []string) (*models.Booking, error)
	CancelBooking(bookingID int64) error
	GetBooking(bookingID int64) (*models.Booking, error)
	GetUserBookings(userID int64) []*models.Booking
}
