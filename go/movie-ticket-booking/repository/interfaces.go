package repository

import "github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"

// MovieRepository defines the interface for movie data operations
type MovieRepository interface {
	AddMovie(movie *models.Movie) error
	GetMovie(movieID int64) (*models.Movie, error)
	SearchMoviesByTitle(title string) []*models.Movie
	SearchMoviesByGenre(genre string) []*models.Movie
	SearchMoviesByLanguage(language string) []*models.Movie
	GetAllMovies() []*models.Movie
}

// TheaterRepository defines the interface for theater data operations
type TheaterRepository interface {
	AddTheater(theater *models.Theater) error
	GetTheater(theaterID int64) (*models.Theater, error)
	GetTheatersByCity(city string) []*models.Theater
	GetAllTheaters() []*models.Theater
}

// ShowRepository defines the interface for show data operations
type ShowRepository interface {
	AddShow(show *models.Show) error
	GetShow(showID int64) (*models.Show, error)
	GetShowsByMovieAndCity(movieID int64, city string) []*models.Show
	GetShowsByTheater(theaterID int64) []*models.Show
	GetShowsByScreen(screenID int64) []*models.Show
	GetAllShows() []*models.Show
}

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	AddBooking(booking *models.Booking) error
	GetBooking(bookingID int64) (*models.Booking, error)
	GetBookingsByUser(userID int64) []*models.Booking
	GetBookingsByShow(showID int64) []*models.Booking
	UpdateBooking(booking *models.Booking) error
}
