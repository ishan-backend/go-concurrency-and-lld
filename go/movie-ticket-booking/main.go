package main

import (
	"fmt"
	"log"

	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/models"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/repository"
	"github.com/ishanpandey/LowLevelDesign/go/movie-ticket-booking/service"
)

func main() {
	// Initialize repositories
	movieRepo := repository.NewMovieRepository()
	theaterRepo := repository.NewTheaterRepository()
	showRepo := repository.NewShowRepository()
	bookingRepo := repository.NewBookingRepository()

	// Initialize services with dependency injection
	movieService := service.NewMovieService(movieRepo)
	theaterService := service.NewTheaterService(theaterRepo)
	showService := service.NewShowService(showRepo, theaterRepo, movieRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, showRepo, theaterRepo)

	fmt.Println("=== Movie Ticket Booking System ===\n")

	// 1. Add Movies
	fmt.Println("=== Adding Movies ===")
	inception, err := movieService.AddMovie("Inception", 148, "Sci-Fi", "English", "2010-07-16")
	if err != nil {
		log.Fatalf("Failed to add movie: %v", err)
	}
	fmt.Printf("Added movie: %s (ID: %d)\n", inception.Title, inception.ID)

	interstellar, err := movieService.AddMovie("Interstellar", 169, "Sci-Fi", "English", "2014-11-07")
	if err != nil {
		log.Fatalf("Failed to add movie: %v", err)
	}
	fmt.Printf("Added movie: %s (ID: %d)\n", interstellar.Title, interstellar.ID)

	avatar, err := movieService.AddMovie("Avatar", 162, "Action", "English", "2009-12-18")
	if err != nil {
		log.Fatalf("Failed to add movie: %v", err)
	}
	fmt.Printf("Added movie: %s (ID: %d)\n\n", avatar.Title, avatar.ID)

	// 2. Add Theaters
	fmt.Println("=== Adding Theaters ===")
	pvrMumbai, err := theaterService.AddTheater("PVR Cinemas", "Mumbai", "Andheri West")
	if err != nil {
		log.Fatalf("Failed to add theater: %v", err)
	}
	fmt.Printf("Added theater: %s (ID: %d) in %s\n", pvrMumbai.Name, pvrMumbai.ID, pvrMumbai.Location.City)

	inoxMumbai, err := theaterService.AddTheater("INOX", "Mumbai", "Goregaon")
	if err != nil {
		log.Fatalf("Failed to add theater: %v", err)
	}
	fmt.Printf("Added theater: %s (ID: %d) in %s\n", inoxMumbai.Name, inoxMumbai.ID, inoxMumbai.Location.City)

	pvrDelhi, err := theaterService.AddTheater("PVR Cinemas", "Delhi", "Connaught Place")
	if err != nil {
		log.Fatalf("Failed to add theater: %v", err)
	}
	fmt.Printf("Added theater: %s (ID: %d) in %s\n\n", pvrDelhi.Name, pvrDelhi.ID, pvrDelhi.Location.City)

	// 3. Add Screens
	fmt.Println("=== Adding Screens ===")
	screen1, err := theaterService.AddScreen(pvrMumbai.ID, 1, 10, 15)
	if err != nil {
		log.Fatalf("Failed to add screen: %v", err)
	}
	fmt.Printf("Added Screen %d to Theater %d: %d rows x %d cols\n", screen1.ScreenNum, pvrMumbai.ID, screen1.TotalRows, screen1.TotalCols)

	screen2, err := theaterService.AddScreen(pvrMumbai.ID, 2, 8, 12)
	if err != nil {
		log.Fatalf("Failed to add screen: %v", err)
	}
	fmt.Printf("Added Screen %d to Theater %d: %d rows x %d cols\n", screen2.ScreenNum, pvrMumbai.ID, screen2.TotalRows, screen2.TotalCols)

	screen3, err := theaterService.AddScreen(inoxMumbai.ID, 1, 12, 18)
	if err != nil {
		log.Fatalf("Failed to add screen: %v", err)
	}
	fmt.Printf("Added Screen %d to Theater %d: %d rows x %d cols\n\n", screen3.ScreenNum, inoxMumbai.ID, screen3.TotalRows, screen3.TotalCols)

	// 4. Create Shows
	fmt.Println("=== Creating Shows ===")
	show1, err := showService.CreateShow(inception.ID, pvrMumbai.ID, screen1.ID, "2024-01-15 18:00", 500.0)
	if err != nil {
		log.Fatalf("Failed to create show: %v", err)
	}
	fmt.Printf("Created show ID: %d for movie %d at %s\n", show1.ID, inception.ID, show1.GetShowTimeOnly())

	show2, err := showService.CreateShow(inception.ID, pvrMumbai.ID, screen1.ID, "2024-01-15 21:00", 600.0)
	if err != nil {
		log.Fatalf("Failed to create show: %v", err)
	}
	fmt.Printf("Created show ID: %d for movie %d at %s\n", show2.ID, inception.ID, show2.GetShowTimeOnly())

	show3, err := showService.CreateShow(interstellar.ID, inoxMumbai.ID, screen3.ID, "2024-01-15 19:30", 550.0)
	if err != nil {
		log.Fatalf("Failed to create show: %v", err)
	}
	fmt.Printf("Created show ID: %d for movie %d at %s\n\n", show3.ID, interstellar.ID, show3.GetShowTimeOnly())

	// 5. Search Movies
	fmt.Println("=== Searching Movies ===")
	sciFiMovies := movieService.SearchMoviesByGenre("Sci-Fi")
	fmt.Printf("Sci-Fi movies: ")
	for _, movie := range sciFiMovies {
		fmt.Printf("%s ", movie.Title)
	}
	fmt.Println()

	englishMovies := movieService.SearchMoviesByLanguage("English")
	fmt.Printf("English movies: ")
	for _, movie := range englishMovies {
		fmt.Printf("%s ", movie.Title)
	}
	fmt.Println("\n")

	// 6. Get Shows by Movie and City
	fmt.Println("=== Getting Shows by Movie and City ===")
	shows := showService.GetShowsByMovieAndCity(inception.ID, "Mumbai")
	fmt.Printf("Shows for '%s' in Mumbai:\n", inception.Title)
	for _, show := range shows {
		fmt.Printf("  Show ID: %d, Time: %s, Price: ₹%.2f\n", show.ID, show.GetShowTimeOnly(), show.Price)
	}
	fmt.Println()

	// 7. View Available Seats
	fmt.Println("=== Viewing Available Seats ===")
	seatLayout, err := showService.GetAvailableSeats(show1.ID)
	if err != nil {
		log.Fatalf("Failed to get seats: %v", err)
	}
	fmt.Printf("Seat layout for Show %d (showing first 5 rows):\n", show1.ID)
	for i := 0; i < 5 && i < len(seatLayout); i++ {
		for j := 0; j < len(seatLayout[i]) && j < 10; j++ {
			seat := seatLayout[i][j]
			status := "A"
			if seat.Status == models.SeatStatusBooked {
				status = "B"
			} else if seat.Status == models.SeatStatusLocked {
				status = "L"
			}
			fmt.Printf("%s[%s] ", seat.ID, status)
		}
		fmt.Println()
	}
	fmt.Println()

	// 8. Create Users
	fmt.Println("=== Creating Users ===")
	user1 := models.NewUser("John Doe", "john@example.com", "1234567890")
	fmt.Printf("Created user: %s (ID: %d)\n", user1.Name, user1.ID)

	user2 := models.NewUser("Jane Smith", "jane@example.com", "0987654321")
	fmt.Printf("Created user: %s (ID: %d)\n\n", user2.Name, user2.ID)

	// 9. Book Tickets
	fmt.Println("=== Booking Tickets ===")
	booking1, err := bookingService.BookTickets(user1.ID, show1.ID, []string{"A1", "A2", "A3"})
	if err != nil {
		log.Fatalf("Failed to book tickets: %v", err)
	}
	fmt.Printf("Booking successful! Booking ID: %d, Total Amount: ₹%.2f\n", booking1.ID, booking1.TotalAmount)
	fmt.Printf("Seats booked: %v\n", booking1.SeatIDs)

	booking2, err := bookingService.BookTickets(user2.ID, show1.ID, []string{"B5", "B6"})
	if err != nil {
		log.Fatalf("Failed to book tickets: %v", err)
	}
	fmt.Printf("Booking successful! Booking ID: %d, Total Amount: ₹%.2f\n", booking2.ID, booking2.TotalAmount)
	fmt.Printf("Seats booked: %v\n\n", booking2.SeatIDs)

	// 10. View Updated Seat Layout
	fmt.Println("=== Updated Seat Layout ===")
	seatLayout, err = showService.GetAvailableSeats(show1.ID)
	if err != nil {
		log.Fatalf("Failed to get seats: %v", err)
	}
	fmt.Printf("Seat layout for Show %d (showing first 5 rows):\n", show1.ID)
	for i := 0; i < 5 && i < len(seatLayout); i++ {
		for j := 0; j < len(seatLayout[i]) && j < 10; j++ {
			seat := seatLayout[i][j]
			status := "A"
			if seat.Status == models.SeatStatusBooked {
				status = "B"
			} else if seat.Status == models.SeatStatusLocked {
				status = "L"
			}
			fmt.Printf("%s[%s] ", seat.ID, status)
		}
		fmt.Println()
	}
	fmt.Println()

	// 11. Get User Bookings
	fmt.Println("=== User Bookings ===")
	user1Bookings := bookingService.GetUserBookings(user1.ID)
	fmt.Printf("Bookings for %s:\n", user1.Name)
	for _, booking := range user1Bookings {
		fmt.Printf("  Booking ID: %d, Show ID: %d, Seats: %v, Amount: ₹%.2f, Status: %s\n",
			booking.ID, booking.ShowID, booking.SeatIDs, booking.TotalAmount, booking.Status)
	}
	fmt.Println()

	// 12. Cancel Booking
	fmt.Println("=== Cancelling Booking ===")
	err = bookingService.CancelBooking(booking1.ID)
	if err != nil {
		log.Fatalf("Failed to cancel booking: %v", err)
	}
	fmt.Printf("Booking %d cancelled successfully\n\n", booking1.ID)

	// 13. View Final Seat Layout
	fmt.Println("=== Final Seat Layout After Cancellation ===")
	seatLayout, err = showService.GetAvailableSeats(show1.ID)
	if err != nil {
		log.Fatalf("Failed to get seats: %v", err)
	}
	fmt.Printf("Seat layout for Show %d (showing first 5 rows):\n", show1.ID)
	for i := 0; i < 5 && i < len(seatLayout); i++ {
		for j := 0; j < len(seatLayout[i]) && j < 10; j++ {
			seat := seatLayout[i][j]
			status := "A"
			if seat.Status == models.SeatStatusBooked {
				status = "B"
			} else if seat.Status == models.SeatStatusLocked {
				status = "L"
			}
			fmt.Printf("%s[%s] ", seat.ID, status)
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
}

