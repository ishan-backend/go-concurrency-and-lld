package service

import (
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/repository"
)

// MovieServiceImpl implements MovieService
type MovieServiceImpl struct {
	movieRepo repository.MovieRepository
}

// NewMovieService creates a new MovieService instance
func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &MovieServiceImpl{
		movieRepo: movieRepo,
	}
}

// AddMovie adds a movie to the system
func (s *MovieServiceImpl) AddMovie(title string, duration int, genre, language string, releaseDate string) (*models.Movie, error) {
	date, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return nil, err
	}
	
	movie := models.NewMovie(title, duration, genre, language, date)
	err = s.movieRepo.AddMovie(movie)
	if err != nil {
		return nil, err
	}
	
	return movie, nil
}

// GetMovie retrieves a movie by ID
func (s *MovieServiceImpl) GetMovie(movieID int64) (*models.Movie, error) {
	return s.movieRepo.GetMovie(movieID)
}

// SearchMoviesByTitle searches movies by title
func (s *MovieServiceImpl) SearchMoviesByTitle(title string) []*models.Movie {
	return s.movieRepo.SearchMoviesByTitle(title)
}

// SearchMoviesByGenre searches movies by genre
func (s *MovieServiceImpl) SearchMoviesByGenre(genre string) []*models.Movie {
	return s.movieRepo.SearchMoviesByGenre(genre)
}

// SearchMoviesByLanguage searches movies by language
func (s *MovieServiceImpl) SearchMoviesByLanguage(language string) []*models.Movie {
	return s.movieRepo.SearchMoviesByLanguage(language)
}

