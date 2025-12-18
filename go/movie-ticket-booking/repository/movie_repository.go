package repository

import (
	"strings"
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
)

// MovieRepositoryImpl implements MovieRepository
type MovieRepositoryImpl struct {
	movies map[int64]*models.Movie
	mu     sync.RWMutex
}

// NewMovieRepository creates a new MovieRepository instance
func NewMovieRepository() MovieRepository {
	return &MovieRepositoryImpl{
		movies: make(map[int64]*models.Movie),
	}
}

// AddMovie adds a movie to the repository
func (r *MovieRepositoryImpl) AddMovie(movie *models.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.movies[movie.ID] = movie
	return nil
}

// GetMovie retrieves a movie by ID
func (r *MovieRepositoryImpl) GetMovie(movieID int64) (*models.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	movie, exists := r.movies[movieID]
	if !exists {
		return nil, models.ErrMovieNotFound
	}
	return movie, nil
}

// SearchMoviesByTitle searches movies by title (case-insensitive partial match)
func (r *MovieRepositoryImpl) SearchMoviesByTitle(title string) []*models.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	title = strings.ToLower(title)
	var results []*models.Movie
	
	for _, movie := range r.movies {
		if strings.Contains(strings.ToLower(movie.Title), title) {
			results = append(results, movie)
		}
	}
	
	return results
}

// SearchMoviesByGenre searches movies by genre (case-insensitive)
func (r *MovieRepositoryImpl) SearchMoviesByGenre(genre string) []*models.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	genre = strings.ToLower(genre)
	var results []*models.Movie
	
	for _, movie := range r.movies {
		if strings.ToLower(movie.Genre) == genre {
			results = append(results, movie)
		}
	}
	
	return results
}

// SearchMoviesByLanguage searches movies by language (case-insensitive)
func (r *MovieRepositoryImpl) SearchMoviesByLanguage(language string) []*models.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	language = strings.ToLower(language)
	var results []*models.Movie
	
	for _, movie := range r.movies {
		if strings.ToLower(movie.Language) == language {
			results = append(results, movie)
		}
	}
	
	return results
}

// GetAllMovies returns all movies
func (r *MovieRepositoryImpl) GetAllMovies() []*models.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	movies := make([]*models.Movie, 0, len(r.movies))
	for _, movie := range r.movies {
		movies = append(movies, movie)
	}
	
	return movies
}

