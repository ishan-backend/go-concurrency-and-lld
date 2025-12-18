package repository

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
)

// TheaterRepositoryImpl implements TheaterRepository
type TheaterRepositoryImpl struct {
	theaters map[int64]*models.Theater
	mu       sync.RWMutex
}

// NewTheaterRepository creates a new TheaterRepository instance
func NewTheaterRepository() TheaterRepository {
	return &TheaterRepositoryImpl{
		theaters: make(map[int64]*models.Theater),
	}
}

// AddTheater adds a theater to the repository
func (r *TheaterRepositoryImpl) AddTheater(theater *models.Theater) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.theaters[theater.ID] = theater
	return nil
}

// GetTheater retrieves a theater by ID
func (r *TheaterRepositoryImpl) GetTheater(theaterID int64) (*models.Theater, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	theater, exists := r.theaters[theaterID]
	if !exists {
		return nil, models.ErrTheaterNotFound
	}
	return theater, nil
}

// GetTheatersByCity retrieves all theaters in a city
func (r *TheaterRepositoryImpl) GetTheatersByCity(city string) []*models.Theater {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var results []*models.Theater
	for _, theater := range r.theaters {
		if theater.Location.City == city {
			results = append(results, theater)
		}
	}
	
	return results
}

// GetAllTheaters returns all theaters
func (r *TheaterRepositoryImpl) GetAllTheaters() []*models.Theater {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	theaters := make([]*models.Theater, 0, len(r.theaters))
	for _, theater := range r.theaters {
		theaters = append(theaters, theater)
	}
	
	return theaters
}

