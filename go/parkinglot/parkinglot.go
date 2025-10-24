// parking_lot_lld.go
// Simplified LLD for a parking lot in Go.
// Features: multi-entry/exit, payments (strategy), multiple vehicle types, multi-floors,
// vacancy "shower" (display), concurrency-safe, tickets.

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// -------------------- Types & Constants --------------------

type VehicleType int

const (
	Motorcycle VehicleType = iota
	Car
	Bus
)

func (vt VehicleType) String() string {
	switch vt {
	case Motorcycle:
		return "Motorcycle"
	case Car:
		return "Car"
	case Bus:
		return "Bus"
	default:
		return "Unknown"
	}
}

// Parking size roughly maps to vehicle type
type SpotSize int

const (
	Small SpotSize = iota
	Medium
	Large
)

// -------------------- Vehicle --------------------

type Vehicle interface {
	GetPlate() string
	GetType() VehicleType
}

type BaseVehicle struct {
	Plate string
	VType VehicleType
}

func (v *BaseVehicle) GetPlate() string { return v.Plate }
func (v *BaseVehicle) GetType() VehicleType { return v.VType }

// Convenience constructors
nfunc NewVehicle(plate string, vt VehicleType) Vehicle {
	return &BaseVehicle{Plate: plate, VType: vt}
}

// -------------------- Parking Spot --------------------

type ParkingSpot struct {
	ID       string
	Size     SpotSize
	Occupied bool
	Vehicle  Vehicle
	mu       sync.Mutex
}

func (ps *ParkingSpot) Assign(v Vehicle) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.Occupied {
		return errors.New("spot already occupied")
	}
	ps.Vehicle = v
	ps.Occupied = true
	return nil
}

func (ps *ParkingSpot) Free() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.Vehicle = nil
	ps.Occupied = false
}

// -------------------- Ticket --------------------

type Ticket struct {
	TicketID    string
	Plate       string
	EntryTime   time.Time
	FloorID     string
	SpotID      string
	EntryGateID string
}

// -------------------- Payment (Strategy Pattern) --------------------

type PaymentResult struct {
	Success bool
	Message string
}

type PaymentStrategy interface {
	Pay(amount float64) PaymentResult
}

// Implementations

type CashPayment struct{}
func (c CashPayment) Pay(amount float64) PaymentResult {
	return PaymentResult{true, fmt.Sprintf("Paid cash: %.2f", amount)}
}

type CardPayment struct{ CardNumber string }
func (c CardPayment) Pay(amount float64) PaymentResult {
	// In real world this would call a gateway
	return PaymentResult{true, fmt.Sprintf("Card %s charged: %.2f", c.CardNumber, amount)}
}

type MobilePayment struct{ Provider string }
func (m MobilePayment) Pay(amount float64) PaymentResult {
	return PaymentResult{true, fmt.Sprintf("Mobile (%s) charged: %.2f", m.Provider, amount)}
}

// -------------------- Observer: Vacancy Display --------------------

// VacancyDisplay listens for floor updates and shows vacancy counts.

type VacancyDisplay struct {
	ID string
	mu sync.Mutex
	// simple map floorID -> vacancies
	Data map[string]int
}

func NewVacancyDisplay(id string) *VacancyDisplay {
	return &VacancyDisplay{ID: id, Data: make(map[string]int)}
}

func (vd *VacancyDisplay) Update(floorID string, vacancies int) {
	vd.mu.Lock()
	defer vd.mu.Unlock()
	vd.Data[floorID] = vacancies
	vd.render()
}

func (vd *VacancyDisplay) render() {
	// In a real system, send to hardware display or web UI
	fmt.Printf("[VacancyDisplay %s] Current vacancies: ")
	for f, v := range vd.Data {
		fmt.Printf("%s=%d ", f, v)
	}
	fmt.Println()
}

// -------------------- Floor --------------------

type Floor struct {
	ID           string
	Spots        []*ParkingSpot
	spotIdxMutex sync.Mutex
	observers    []*VacancyDisplay
}

func NewFloor(id string, counts map[SpotSize]int) *Floor {
	f := &Floor{ID: id}
	for size, cnt := range counts {
		for i := 0; i < cnt; i++ {
			ps := &ParkingSpot{ID: fmt.Sprintf("%s-%d", id, len(f.Spots)), Size: size}
			f.Spots = append(f.Spots, ps)
		}
	}
	return f
}

func (f *Floor) RegisterObserver(vd *VacancyDisplay) {
	f.observers = append(f.observers, vd)
}

func (f *Floor) notify() {
	vac := f.Vacancies()
	for _, ob := range f.observers {
		ob.Update(f.ID, vac)
	}
}

func (f *Floor) Vacancies() int {
	count := 0
	for _, s := range f.Spots {
		if !s.Occupied {
			count++
		}
	}
	return count
}

// find spot for vehicle (simple first-fit)
func (f *Floor) FindAndAssign(v Vehicle) (*ParkingSpot, error) {
	f.spotIdxMutex.Lock()
	defer f.spotIdxMutex.Unlock()
	for _, s := range f.Spots {
		if !s.Occupied && fits(v.GetType(), s.Size) {
			if err := s.Assign(v); err == nil {
				f.notify()
				return s, nil
			}
		}
	}
	return nil, errors.New("no spot available on this floor")
}

func (f *Floor) ReleaseSpot(spotID string) bool {
	for _, s := range f.Spots {
		if s.ID == spotID {
			s.Free()
			f.notify()
			return true
		}
	}
	return false
}

func fits(vt VehicleType, sz SpotSize) bool {
	switch vt {
	case Motorcycle:
		return true // motorcycles fit anywhere
	case Car:
		return sz == Medium || sz == Large
	case Bus:
		return sz == Large
	}
	return false
}

// -------------------- Gate/Entry/Exit --------------------

type Gate struct {
	ID string
}

// -------------------- ParkingLot (Singleton) --------------------

type ParkingLot struct {
	Floors      map[string]*Floor
	Entries     map[string]*Gate
	Exits       map[string]*Gate
	Tickets     map[string]*Ticket
	muTickets   sync.RWMutex
	muFloors    sync.RWMutex
	rand        *rand.Rand
}

var (
	lotInstance *ParkingLot
	once        sync.Once
)

func GetParkingLot() *ParkingLot {
	once.Do(func() {
		lotInstance = &ParkingLot{
			Floors:  make(map[string]*Floor),
			Entries: make(map[string]*Gate),
			Exits:   make(map[string]*Gate),
			Tickets: make(map[string]*Ticket),
			rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		}
	})
	return lotInstance
}

func (pl *ParkingLot) AddFloor(f *Floor) {
	pl.muFloors.Lock()
	defer pl.muFloors.Unlock()
	pl.Floors[f.ID] = f
}

func (pl *ParkingLot) AddEntry(g *Gate) { pl.Entries[g.ID] = g }
func (pl *ParkingLot) AddExit(g *Gate) { pl.Exits[g.ID] = g }

// Entry: issues ticket and assigns spot (search floors in order)
func (pl *ParkingLot) Enter(plate string, vt VehicleType, entryGateID string) (*Ticket, error) {
	v := &BaseVehicle{Plate: plate, VType: vt}
	pl.muFloors.RLock()
	defer pl.muFloors.RUnlock()
	for _, f := range pl.Floors {
		if spot, err := f.FindAndAssign(v); err == nil {
			// create ticket
			id := fmt.Sprintf("T-%d", pl.rand.Intn(1e9))
			t := &Ticket{TicketID: id, Plate: plate, EntryTime: time.Now(), FloorID: f.ID, SpotID: spot.ID, EntryGateID: entryGateID}
			pl.muTickets.Lock()
			pl.Tickets[id] = t
			pl.muTickets.Unlock()
			return t, nil
		}
	}
	return nil, errors.New("parking full")
}

// Exit: calculate fee -> use payment strategy, free spot
func (pl *ParkingLot) Exit(ticketID string, payment PaymentStrategy) (PaymentResult, error) {
	pl.muTickets.RLock()
	t, ok := pl.Tickets[ticketID]
	pl.muTickets.RUnlock()
	if !ok {
		return PaymentResult{false, "invalid ticket"}, errors.New("invalid ticket")
	}
	// simple fee calculation
	dur := time.Since(t.EntryTime)
	hours := dur.Hours()
	if hours < 1 { hours = 1 }
	amount := 10.0 * hours // Rs 10 per hour (example)
	res := payment.Pay(amount)
	if res.Success {
		// free spot
		pl.muFloors.RLock()
		floor, ok := pl.Floors[t.FloorID]
		pl.muFloors.RUnlock()
		if ok {
			floor.ReleaseSpot(t.SpotID)
		}
		pl.muTickets.Lock()
		delete(pl.Tickets, ticketID)
		pl.muTickets.Unlock()
	}
	return res, nil
}

// Quick status helper
func (pl *ParkingLot) Status() {
	pl.muFloors.RLock()
	defer pl.muFloors.RUnlock()
	for id, f := range pl.Floors {
		fmt.Printf("Floor %s vacancies=%d\n", id, f.Vacancies())
	}
}

// -------------------- Example usage --------------------

func main() {
	pl := GetParkingLot()

	// Create vacancy display and floors
	vd := NewVacancyDisplay("Main-Display")

	f1 := NewFloor("F1", map[SpotSize]int{Small: 2, Medium: 4, Large: 1})
	f2 := NewFloor("F2", map[SpotSize]int{Small: 4, Medium: 6, Large: 2})

	f1.RegisterObserver(vd)
	f2.RegisterObserver(vd)

	pl.AddFloor(f1)
	pl.AddFloor(f2)

	pl.AddEntry(&Gate{ID: "E1"})
	pl.AddEntry(&Gate{ID: "E2"})
	pl.AddExit(&Gate{ID: "X1"})

	// Vehicles enter
	t1, err := pl.Enter("MH12AB1111", Car, "E1")
	if err != nil { fmt.Println("enter err:", err) } else { fmt.Println("Ticket issued:", t1) }

	t2, _ := pl.Enter("MH12AB2222", Motorcycle, "E2")
	fmt.Println("Ticket issued:", t2)

	pl.Status()

	// simulate wait
	time.Sleep(2 * time.Second)

	// Exit t1 with card
	res, err := pl.Exit(t1.TicketID, CardPayment{CardNumber: "4111-xxxx"})
	fmt.Println("Exit result:", res, err)

	pl.Status()
}

/*
Design Patterns used (short comments):
- Singleton: ParkingLot is a singleton to provide a single source of truth and coordinated state.
- Factory (light): NewFloor/NewVehicle helpers act as simple factories for creating objects.
- Strategy: PaymentStrategy lets you swap payment methods (cash/card/mobile).
- Observer: VacancyDisplay registers to Floors to receive vacancy updates.
- Mutex/Concurrency: not a pattern but essential for thread-safety in LLD.
- Template Method / Command (can be applied): fee calculation and payment handling could be extended with templating or command objects for complex flows.
- Composite: Floor contains multiple ParkingSpots (tree-like) — can be extended for hierarchical queries.
*/
