# Movie Ticket Booking System - Low Level Design

## Problem Statement

Design and implement a Movie Ticket Booking System that allows users to search for movies, view available shows, select seats, and book tickets. The system should handle seat availability, booking management, and support multiple theaters with multiple screens.

### Core Requirements

#### 1. Movie Management
- **Add Movie**: Add movies with details:
  - Title
  - Duration (in minutes)
  - Genre
  - Language
  - Release Date
- **Search Movies**: Search movies by title, genre, or language
- **Get Movie**: Retrieve movie details by ID

#### 2. Theater Management
- **Add Theater**: Add theaters with:
  - Name
  - Location (City, Address)
  - List of Screens
- **Add Screen**: Add screens to theaters with:
  - Screen Number
  - Total Seats (rows and columns)
  - Seat Layout (Normal, Premium, VIP)
- **Get Theaters**: Get theaters by location/city

#### 3. Show Management
- **Create Show**: Create shows for movies in theaters:
  - Movie ID
  - Theater ID
  - Screen ID
  - Show Time
  - Show Date
- **Get Shows**: Get available shows for a movie in a city
- **Get Show Details**: Get show details including available seats

#### 4. Seat Management
- **View Seats**: View seat layout for a show
- **Seat Types**: Support different seat types (Normal, Premium, VIP)
- **Seat Status**: Track seat status (Available, Booked, Locked)

#### 5. Booking Management
- **Book Tickets**: Book tickets for a show:
  - User ID
  - Show ID
  - List of Seat IDs
- **Cancel Booking**: Cancel existing bookings
- **Get Booking**: Retrieve booking details
- **Get User Bookings**: Get all bookings for a user

### Technical Constraints

1. **Concurrency**: Handle concurrent bookings for the same seats
2. **Seat Locking**: Lock seats temporarily during booking process (e.g., 5 minutes)
3. **Data Structures**: Use efficient data structures for fast lookups
4. **Thread Safety**: Ensure thread-safe operations for seat booking

### Expected Functionalities

```
1. Add movies to the system
2. Add theaters with screens
3. Create shows for movies in theaters
4. Search movies by various criteria
5. Get available shows for a movie in a city
6. View seat layout for a show
7. Book tickets (with seat locking mechanism)
8. Cancel bookings
9. View booking history
```

### Example Usage Flow

```
1. Add Movie: "Inception", Duration: 148 mins, Genre: "Sci-Fi", Language: "English"
2. Add Theater: "PVR Cinemas", Location: "Mumbai"
3. Add Screen: Screen 1 with 10 rows, 15 columns
4. Create Show: Movie "Inception" at PVR Cinemas, Screen 1, Time: 18:00, Date: 2024-01-15
5. Search Movies: Search by genre "Sci-Fi"
6. Get Shows: Get all shows for "Inception" in "Mumbai"
7. View Seats: View available seats for a show
8. Book Tickets: Book seats A1, A2, A3 for user
9. Cancel Booking: Cancel the booking
10. Get User Bookings: View all bookings for a user
```

---

## SOLID Principles Application

### 1. Single Responsibility Principle (SRP)
- **MovieService**: Handles only movie-related operations
- **TheaterService**: Handles only theater and screen operations
- **ShowService**: Handles only show-related operations
- **BookingService**: Handles only booking operations
- Each repository handles data access for a specific entity

### 2. Open/Closed Principle (OCP)
- Use interfaces for services and repositories
- Allow extension without modification (e.g., different payment methods, seat types)
- Strategy pattern for different booking strategies

### 3. Liskov Substitution Principle (LSP)
- Repository implementations are interchangeable
- Service implementations follow the same contract

### 4. Interface Segregation Principle (ISP)
- Separate interfaces for:
  - `MovieRepository` (movie data operations)
  - `TheaterRepository` (theater data operations)
  - `ShowRepository` (show data operations)
  - `BookingRepository` (booking data operations)

### 5. Dependency Inversion Principle (DIP)
- Services depend on repository interfaces, not concrete implementations
- Use dependency injection

---

## Design Considerations

### Data Structures

1. **Movie Storage**:
   - Map for O(1) lookup by ID
   - Index by genre, language for efficient search

2. **Theater Storage**:
   - Map for O(1) lookup by ID
   - Index by city for location-based queries

3. **Show Storage**:
   - Map for O(1) lookup by ID
   - Index by movie ID and city for efficient queries

4. **Seat Storage**:
   - 2D array/matrix for seat layout
   - Map for O(1) seat lookup by ID
   - Track seat status and type

5. **Booking Storage**:
   - Map for O(1) lookup by booking ID
   - Index by user ID for user bookings
   - Index by show ID for show bookings

### Seat Locking Mechanism

- Implement a temporary lock on seats during booking
- Lock expires after a timeout (e.g., 5 minutes)
- Prevent double booking of locked seats

### Architecture Layers

```
┌─────────────────────────────────────┐
│         Service Layer               │
│  (MovieService, TheaterService,     │
│   ShowService, BookingService)      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Repository Layer                │
│  (MovieRepository, TheaterRepository│
│   ShowRepository, BookingRepository)│
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Model Layer                  │
│  (Movie, Theater, Screen, Show,     │
│   Seat, Booking, User)              │
└─────────────────────────────────────┘
```

### Key Design Patterns

1. **Repository Pattern**: Abstract data access layer
2. **Service Pattern**: Business logic encapsulation
3. **Factory Pattern**: For creating service instances
4. **Strategy Pattern**: For different booking strategies (optional)
5. **Observer Pattern**: For notifications (optional)

---

## Implementation Guidelines

### File Structure

```
movie-ticket-booking/
├── README.md
├── go.mod
├── models/
│   ├── movie.go
│   ├── theater.go
│   ├── screen.go
│   ├── show.go
│   ├── seat.go
│   ├── booking.go
│   └── user.go
├── repository/
│   ├── movie_repository.go
│   ├── theater_repository.go
│   ├── show_repository.go
│   ├── booking_repository.go
│   └── interfaces.go
├── service/
│   ├── movie_service.go
│   ├── theater_service.go
│   ├── show_service.go
│   ├── booking_service.go
│   └── interfaces.go
└── main.go
```

---

## Future Enhancements (Optional)

1. **Payment Integration**: Integrate payment gateway
2. **Notifications**: Email/SMS notifications for bookings
3. **Waitlist**: Waitlist for fully booked shows
4. **Discounts**: Support for discounts and coupons
5. **Reviews**: Movie and theater reviews
6. **Recommendations**: Movie recommendations based on user history
7. **Multiple Bookings**: Support for booking multiple shows at once
8. **Refund Policy**: Automated refund processing

