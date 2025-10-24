// Simplified Splitwise LLD in Go
// Version: uses plain string IDs instead of custom ID type

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// ------------------------- Models -------------------------

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type Expense struct {
	ID        string             `json:"id"`
	GroupID   *string            `json:"group_id,omitempty"`
	Title     string             `json:"title"`
	Amount    float64            `json:"amount"`
	Payer     string             `json:"payer"`
	CreatedBy string             `json:"created_by"`
	CreatedAt time.Time          `json:"created_at"`
	Shares    map[string]float64 `json:"shares"`
}

type Transaction struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// ------------------------- Store -------------------------

type Store struct {
	mu           sync.RWMutex
	users        map[string]*User
	groups       map[string]*Group
	expenses     map[string]*Expense
	transactions []*Transaction
	balances     map[string]map[string]float64 // user -> (user -> amount)
}

func NewStore() *Store {
	return &Store{
		users:    make(map[string]*User),
		groups:   make(map[string]*Group),
		expenses: make(map[string]*Expense),
		balances: make(map[string]map[string]float64),
	}
}

// ------------------------- Helpers -------------------------

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func (s *Store) addUser(u *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
	if _, ok := s.balances[u.ID]; !ok {
		s.balances[u.ID] = make(map[string]float64)
	}
}

func (s *Store) addGroup(g *Group) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups[g.ID] = g
	for _, m := range g.Members {
		if _, ok := s.balances[m]; !ok {
			s.balances[m] = make(map[string]float64)
		}
	}
}

func (s *Store) recordExpense(e *Expense) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expenses[e.ID] = e
	for uid, amt := range e.Shares {
		if uid == e.Payer {
			continue
		}
		s.balances[uid][e.Payer] += amt
		s.balances[e.Payer][uid] -= amt
	}
}

func (s *Store) settle(from, to string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.balances[from][to] -= amount
	s.balances[to][from] += amount

	tr := &Transaction{
		ID:        newID("txn"),
		From:      from,
		To:        to,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	s.transactions = append(s.transactions, tr)
	return tr, nil
}

func (s *Store) getNetBalances() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string]float64)
	for from, m := range s.balances {
		for to, amt := range m {
			res[from] -= amt
			res[to] += amt
		}
	}
	return res
}

// ------------------------- Service -------------------------

type Service struct {
	store *Store
}

func NewService(st *Store) *Service {
	return &Service{store: st}
}

func (svc *Service) CreateUser(name string) *User {
	u := &User{ID: newID("user"), Name: name}
	svc.store.addUser(u)
	return u
}

func (svc *Service) CreateGroup(name string, members []string) *Group {
	g := &Group{ID: newID("grp"), Name: name, Members: members}
	svc.store.addGroup(g)
	return g
}

func (svc *Service) AddExpense(title string, amount float64, payer, createdBy string, groupID *string, shares map[string]float64) (*Expense, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	sum := 0.0
	for _, v := range shares {
		sum += v
	}
	if math.Abs(sum-amount) > 0.01 {
		return nil, fmt.Errorf("shares don't sum up to amount")
	}
	e := &Expense{
		ID:        newID("exp"),
		GroupID:   groupID,
		Title:     title,
		Amount:    amount,
		Payer:     payer,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		Shares:    shares,
	}
	svc.store.recordExpense(e)
	return e, nil
}

func (svc *Service) AutoSettleAll() ([]*Transaction, error) {
	net := svc.store.getNetBalances()
	creditors := []struct {
		uid string
		amt float64
	}{}
	debtors := []struct {
		uid string
		amt float64
	}{}
	for uid, b := range net {
		if math.Abs(b) < 0.01 {
			continue
		}
		if b > 0 {
			creditors = append(creditors, struct {
				uid string
				amt float64
			}{uid, b})
		} else {
			debtors = append(debtors, struct {
				uid string
				amt float64
			}{uid, -b}) // +ve owed amount
		}
	}
	sort.Slice(creditors, func(i, j int) bool { return creditors[i].amt > creditors[j].amt }) // biggest creditor ahead
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].amt > debtors[j].amt })

	var results []*Transaction
	ci, di := 0, 0
	for ci < len(creditors) && di < len(debtors) {
		c := creditors[ci]
		d := debtors[di]
		minAmt := math.Min(c.amt, d.amt) // min amount which can be settled
		tr, err := svc.store.settle(d.uid, c.uid, minAmt)
		if err != nil {
			return results, err
		}
		results = append(results, tr)
		creditors[ci].amt -= minAmt
		debtors[di].amt -= minAmt
		if creditors[ci].amt <= 0.01 {
			ci++
		}
		if debtors[di].amt <= 0.01 {
			di++
		}
	}
	return results, nil
}

/*
	•	In worst case, O(N!) — but for small groups (N ≤ 8), it’s fine.
	•	It’s optimal in terms of minimizing transactions.
	•	Production Splitwise usually uses a greedy (two-pointer) version because it scales better.

func minTransfers(balances []int) int {
    return dfs(0, balances)
}

func dfs(start int, balances []int) int {
    // skip settled users
    for start < len(balances) && balances[start] == 0 {
        start++
    }
    if start == len(balances) {
        return 0
    }

    res := math.MaxInt
    for i := start + 1; i < len(balances); i++ {
        // opposite signs → can settle
        if balances[start]*balances[i] < 0 {
            // try to settle start with i
            balances[i] += balances[start]
            res = min(res, 1+dfs(start+1, balances))
            balances[i] -= balances[start] // backtrack
        }
    }
    return res
}
*/

// ------------------------- HTTP -------------------------

var svc *Service

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSON(w, map[string]string{"error": "name required"}, 400)
		return
	}
	u := svc.CreateUser(name)
	writeJSON(w, u, 201)
}

func handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	g := svc.CreateGroup(req.Name, req.Members)
	writeJSON(w, g, 201)
}

func handleAddExpense(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title     string             `json:"title"`
		Amount    float64            `json:"amount"`
		Payer     string             `json:"payer"`
		CreatedBy string             `json:"created_by"`
		GroupID   *string            `json:"group_id,omitempty"`
		Shares    map[string]float64 `json:"shares"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	e, err := svc.AddExpense(req.Title, req.Amount, req.Payer, req.CreatedBy, req.GroupID, req.Shares)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	writeJSON(w, e, 201)
}

func handleGetNetBalances(w http.ResponseWriter, r *http.Request) {
	res := svc.store.getNetBalances()
	writeJSON(w, res, 200)
}

func handleAutoSettle(w http.ResponseWriter, r *http.Request) {
	trs, err := svc.AutoSettleAll()
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 500)
		return
	}
	writeJSON(w, trs, 200)
}

func handleManualSettle(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amtStr := r.URL.Query().Get("amount")
	amt, err := strconv.ParseFloat(amtStr, 64)
	if from == "" || to == "" || err != nil {
		writeJSON(w, map[string]string{"error": "invalid params"}, 400)
		return
	}
	tr, err := svc.store.settle(from, to, amt)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	writeJSON(w, tr, 200)
}

func main() {
	store := NewStore()
	svc = NewService(store)

	a := svc.CreateUser("Alice")
	b := svc.CreateUser("Bob")
	c := svc.CreateUser("Charlie")
	g := svc.CreateGroup("Trip", []string{a.ID, b.ID, c.ID})

	shares := map[string]float64{a.ID: 100, b.ID: 100, c.ID: 100}
	svc.AddExpense("Dinner", 300, a.ID, a.ID, &g.ID, shares)

	http.HandleFunc("/users/create", handleCreateUser)
	http.HandleFunc("/groups/create", handleCreateGroup)
	http.HandleFunc("/expenses/add", handleAddExpense)
	http.HandleFunc("/balances/net", handleGetNetBalances)
	http.HandleFunc("/settle/auto", handleAutoSettle)
	http.HandleFunc("/settle/manual", handleManualSettle)

	log.Println("Running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
