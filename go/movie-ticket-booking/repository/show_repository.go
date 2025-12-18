package repository

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
)

// ShowRepositoryImpl implements ShowRepository
type ShowRepositoryImpl struct {
	shows map[int64]*models.Show
	mu    sync.RWMutex
}

// NewShowRepository creates a new ShowRepository instance
func NewShowRepository() ShowRepository {
	return &ShowRepositoryImpl{
		shows: make(map[int64]*models.Show),
	}
}

// AddShow adds a show to the repository
func (r *ShowRepositoryImpl) AddShow(show *models.Show) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shows[show.ID] = show
	return nil
}

// GetShow retrieves a show by ID
func (r *ShowRepositoryImpl) GetShow(showID int64) (*models.Show, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	show, exists := r.shows[showID]
	if !exists {
		return nil, models.ErrShowNotFound
	}
	return show, nil
}

// GetShowsByMovieAndCity retrieves shows for a movie in a specific city
func (r *ShowRepositoryImpl) GetShowsByMovieAndCity(movieID int64, city string) []*models.Show {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Show
	for _, show := range r.shows {
		if show.MovieID == movieID {
			// Note: We need theater info to check city, but we don't have direct access here
			// This will be handled in the service layer
			results = append(results, show)
		}
	}

	return results
}

// GetShowsByTheater retrieves all shows for a theater
func (r *ShowRepositoryImpl) GetShowsByTheater(theaterID int64) []*models.Show {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Show
	for _, show := range r.shows {
		if show.TheaterID == theaterID {
			results = append(results, show)
		}
	}

	return results
}

// GetShowsByScreen retrieves all shows for a screen
func (r *ShowRepositoryImpl) GetShowsByScreen(screenID int64) []*models.Show {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Show
	for _, show := range r.shows {
		if show.ScreenID == screenID {
			results = append(results, show)
		}
	}

	return results
}

// GetAllShows returns all shows
func (r *ShowRepositoryImpl) GetAllShows() []*models.Show {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shows := make([]*models.Show, 0, len(r.shows))
	for _, show := range r.shows {
		shows = append(shows, show)
	}

	return shows
}
