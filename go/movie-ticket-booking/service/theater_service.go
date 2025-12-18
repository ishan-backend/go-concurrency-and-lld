package service

import (
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/repository"
)

// TheaterServiceImpl implements TheaterService
type TheaterServiceImpl struct {
	theaterRepo repository.TheaterRepository
}

// NewTheaterService creates a new TheaterService instance
func NewTheaterService(theaterRepo repository.TheaterRepository) TheaterService {
	return &TheaterServiceImpl{
		theaterRepo: theaterRepo,
	}
}

// AddTheater adds a theater to the system
func (s *TheaterServiceImpl) AddTheater(name, city, address string) (*models.Theater, error) {
	theater := models.NewTheater(name, city, address)
	err := s.theaterRepo.AddTheater(theater)
	if err != nil {
		return nil, err
	}

	return theater, nil
}

// GetTheater retrieves a theater by ID
func (s *TheaterServiceImpl) GetTheater(theaterID int64) (*models.Theater, error) {
	return s.theaterRepo.GetTheater(theaterID)
}

// AddScreen adds a screen to a theater
func (s *TheaterServiceImpl) AddScreen(theaterID int64, screenNum, totalRows, totalCols int) (*models.Screen, error) {
	theater, err := s.theaterRepo.GetTheater(theaterID)
	if err != nil {
		return nil, err
	}

	screen := models.NewScreen(theaterID, screenNum, totalRows, totalCols)
	theater.AddScreen(screen)

	return screen, nil
}

// GetTheatersByCity retrieves theaters by city
func (s *TheaterServiceImpl) GetTheatersByCity(city string) []*models.Theater {
	return s.theaterRepo.GetTheatersByCity(city)
}
