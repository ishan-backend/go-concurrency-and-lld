package main

import (
	"errors"
	"fmt"
	"sync"
)

// Product represents a product in the vending machine
type Product struct {
	ID    string
	Name  string
	Price int // price in cents
}

// Slot holds product inventory and capacity
type Slot struct {
	Product  Product
	Quantity int
	Capacity int
}

// Inventory using Singleton-like struct (but kept simple and thread-safe)
type Inventory struct {
	slots map[string]*Slot
	mu    sync.RWMutex
}

// NewInventory returns initialized inventory
func NewInventory() *Inventory {
	return &Inventory{slots: make(map[string]*Slot)}
}

func (inv *Inventory) AddSlot(id string, p Product, capacity int) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.slots[id] = &Slot{Product: p, Quantity: 0, Capacity: capacity}
}

func (inv *Inventory) Restock(id string, qty int) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	s, ok := inv.slots[id]
	if !ok {
		return errors.New("slot not found")
	}
	if s.Quantity+qty > s.Capacity {
		return errors.New("exceeds capacity")
	}
	s.Quantity += qty
	return nil
}

func (inv *Inventory) Dispense(id string) (Product, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	s, ok := inv.slots[id]
	if !ok || s.Quantity == 0 {
		return Product{}, errors.New("product unavailable")
	}
	s.Quantity--
	return s.Product, nil
}

func (inv *Inventory) GetSlot(id string) (*Slot, bool) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	s, ok := inv.slots[id]
	return s, ok
}

// Payment strategy interface - allows plugging different payment processors
type PaymentProcessor interface {
	Accept(amount int) error
	Refund() (int, error)
	Balance() int
}

// SimpleCoinProcessor - example implementation that accepts coins / notes
type SimpleCoinProcessor struct {
	mu      sync.Mutex
	balance int
}

func NewSimpleCoinProcessor() *SimpleCoinProcessor { return &SimpleCoinProcessor{} }

func (p *SimpleCoinProcessor) Accept(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.balance += amount
	return nil
}

func (p *SimpleCoinProcessor) Refund() (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	amt := p.balance
	p.balance = 0
	return amt, nil
}

func (p *SimpleCoinProcessor) Balance() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.balance
}

// ChangeMaker - simple greedy change maker for typical coin sets
func MakeChange(amount int, coins []int) (map[int]int, error) {
	if amount < 0 {
		return nil, errors.New("invalid amount")
	}
	res := make(map[int]int)
	for _, c := range coins {
		if c <= 0 {
			continue
		}
		if amount == 0 {
			break
		}
		cnt := amount / c
		if cnt > 0 {
			res[c] = cnt
			amount -= cnt * c
		}
	}
	if amount != 0 {
		return nil, errors.New("cannot provide exact change")
	}
	return res, nil
}

// VendingMachine acts as the facade for operations
type VendingMachine struct {
	inventory     *Inventory
	payment       PaymentProcessor
	acceptedCoins []int
	mu            sync.Mutex
	state         State
}

func NewVendingMachine(inv *Inventory, payment PaymentProcessor, acceptedCoins []int) *VendingMachine {
	vm := &VendingMachine{inventory: inv, payment: payment, acceptedCoins: acceptedCoins}
	vm.state = &IdleState{vm: vm}
	return vm
}

// High level operations
func (vm *VendingMachine) SelectProduct(slotID string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.state.Select(slotID)
}

func (vm *VendingMachine) InsertMoney(amount int) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.state.Insert(amount)
}

func (vm *VendingMachine) Dispense() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.state.Dispense()
}

func (vm *VendingMachine) Cancel() (int, error) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.state.Cancel()
}

// State pattern to manage transitions and responsibilities
type State interface {
	Select(slotID string) error
	Insert(amount int) error
	Dispense() error
	Cancel() (int, error)
}

// IdleState - waiting for user
type IdleState struct{ vm *VendingMachine }

func (s *IdleState) Select(slotID string) error {
	slot, ok := s.vm.inventory.GetSlot(slotID)
	if !ok || slot.Quantity == 0 {
		return errors.New("product not available")
	}
	// store selected product price in a simple way: replace payment processor (or hold in vm)
	s.vm.state = &WaitingForMoneyState{vm: s.vm, selectedSlot: slotID}
	return nil
}

func (s *IdleState) Insert(amount int) error { return errors.New("select product first") }
func (s *IdleState) Dispense() error         { return errors.New("no product selected") }
func (s *IdleState) Cancel() (int, error)    { return 0, errors.New("nothing to cancel") }

// WaitingForMoneyState - product selected, waiting for money
type WaitingForMoneyState struct {
	vm           *VendingMachine
	selectedSlot string
}

func (s *WaitingForMoneyState) Select(slotID string) error {
	return errors.New("product already selected")
}

func (s *WaitingForMoneyState) Insert(amount int) error {
	// basic validation of coin denomination
	valid := false
	for _, c := range s.vm.acceptedCoins {
		if c == amount {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("coin not accepted")
	}
	if err := s.vm.payment.Accept(amount); err != nil {
		return err
	}
	// check if enough money
	slot, _ := s.vm.inventory.GetSlot(s.selectedSlot)
	if s.vm.payment.Balance() >= slot.Product.Price {
		s.vm.state = &DispensingState{vm: s.vm, selectedSlot: s.selectedSlot}
	}
	return nil
}

func (s *WaitingForMoneyState) Dispense() error { return errors.New("insufficient funds") }

func (s *WaitingForMoneyState) Cancel() (int, error) {
	amt, err := s.vm.payment.Refund()
	if err != nil {
		return 0, err
	}
	s.vm.state = &IdleState{vm: s.vm}
	return amt, nil
}

// DispensingState - ready to dispense, calculate change and send product
type DispensingState struct {
	vm           *VendingMachine
	selectedSlot string
}

func (s *DispensingState) Select(slotID string) error { return errors.New("already processing") }

func (s *DispensingState) Insert(amount int) error { return errors.New("already processing") }

func (s *DispensingState) Dispense() error {
	// dispense product
	product, err := s.vm.inventory.Dispense(s.selectedSlot)
	if err != nil {
		s.vm.state = &IdleState{vm: s.vm}
		return err
	}
	paid := s.vm.payment.Balance()
	change := paid - product.Price
	// attempt to make change
	coinsMap, err := MakeChange(change, s.vm.acceptedCoins)
	if err != nil {
		// if cannot make change, rollback: refund full amount and restock product
		_ = s.vm.inventory.Restock(s.selectedSlot, 1)
		amt, _ := s.vm.payment.Refund()
		s.vm.state = &IdleState{vm: s.vm}
		return fmt.Errorf("cannot dispense exact change, refunded %d", amt)
	}
	// success: clear payment balance and transition to idle
	_, _ = s.vm.payment.Refund() // flush balance
	fmt.Printf("Dispensed: %s. Change given: %v\n", product.Name, coinsMap)
	s.vm.state = &IdleState{vm: s.vm}
	return nil
}

func (s *DispensingState) Cancel() (int, error) {
	return 0, errors.New("cannot cancel while dispensing")
}

// Example usage in main
func main() {
	inv := NewInventory()
	inv.AddSlot("A1", Product{ID: "p1", Name: "Coke", Price: 125}, 10)
	inv.AddSlot("B1", Product{ID: "p2", Name: "Pepsi", Price: 100}, 10)
	_ = inv.Restock("A1", 5)
	_ = inv.Restock("B1", 3)

	payment := NewSimpleCoinProcessor()
	coins := []int{100, 50, 25, 10, 5, 1} // accepted denominations in cents
	vm := NewVendingMachine(inv, payment, coins)

	// user flow
	if err := vm.SelectProduct("A1"); err != nil {
		fmt.Println("select error:", err)
		return
	}
	if err := vm.InsertMoney(100); err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.InsertMoney(25); err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.Dispense(); err != nil {
		fmt.Println("dispense err:", err)
		return
	}

	// cancel flow
	if err := vm.SelectProduct("B1"); err != nil {
		fmt.Println("select error:", err)
		return
	}
	if err := vm.InsertMoney(25); err != nil {
		fmt.Println(err)
		return
	}
	amt, _ := vm.Cancel()
	fmt.Println("Refunded:", amt)
}
