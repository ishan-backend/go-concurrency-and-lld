package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var svc *Service

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, message string, code int) {
	writeJSON(w, map[string]string{"error": message}, code)
}

// ------------------------- HTTP Handlers -------------------------

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var name, email string
	if r.Method == http.MethodPost {
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = req.Name
		email = req.Email
	} else {
		name = r.URL.Query().Get("name")
		email = r.URL.Query().Get("email")
	}

	if name == "" {
		writeError(w, "name required", http.StatusBadRequest)
		return
	}
	if email == "" {
		email = fmt.Sprintf("%s@example.com", name)
	}

	u := svc.CreateUser(name, email)
	writeJSON(w, u, http.StatusCreated)
}

func handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
		Admin   string   `json:"admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		writeError(w, "name required", http.StatusBadRequest)
		return
	}
	if len(req.Members) == 0 {
		writeError(w, "members required", http.StatusBadRequest)
		return
	}
	if req.Admin == "" {
		req.Admin = req.Members[0] // Default to first member
	}

	g, err := svc.CreateGroup(req.Name, req.Members, req.Admin)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, g, http.StatusCreated)
}

func handleAddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeError(w, "title required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		writeError(w, "amount must be positive", http.StatusBadRequest)
		return
	}
	if req.Payer == "" {
		writeError(w, "payer required", http.StatusBadRequest)
		return
	}
	if req.CreatedBy == "" {
		req.CreatedBy = req.Payer
	}

	e, err := svc.AddExpense(&req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, e, http.StatusCreated)
}

func handleGetBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user_id from path: /api/users/{user_id}/balances
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "balances" {
		writeError(w, "Invalid path. Expected: /api/users/{user_id}/balances", http.StatusBadRequest)
		return
	}
	userID := parts[len(parts)-2]

	var groupID *string
	if gid := r.URL.Query().Get("group_id"); gid != "" {
		groupID = &gid
	}

	balances, err := svc.GetBalances(userID, groupID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, balances, http.StatusOK)
}

func handleSimplifyDebts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract group_id from path: /api/groups/{group_id}/simplify-debts
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "simplify-debts" {
		writeError(w, "Invalid path. Expected: /api/groups/{group_id}/simplify-debts", http.StatusBadRequest)
		return
	}
	groupIDStr := parts[len(parts)-2]

	var groupID *string
	if groupIDStr != "" && groupIDStr != "global" {
		groupID = &groupIDStr
	}

	result, err := svc.SimplifyDebts(groupID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, result, http.StatusOK)
}

func handleGetActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract group_id from path: /api/groups/{group_id}/activity
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "activity" {
		writeError(w, "Invalid path. Expected: /api/groups/{group_id}/activity", http.StatusBadRequest)
		return
	}
	groupIDStr := parts[len(parts)-2]

	var groupID *string
	if groupIDStr != "" && groupIDStr != "global" {
		groupID = &groupIDStr
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	activities, err := svc.GetActivities(groupID, limit, offset)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, activities, http.StatusOK)
}

func handleSetDefaultSplit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract group_id from path: /api/groups/{group_id}/default-split
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[len(parts)-1] != "default-split" {
		writeError(w, "Invalid path. Expected: /api/groups/{group_id}/default-split", http.StatusBadRequest)
		return
	}
	groupID := parts[len(parts)-2]

	var req SetDefaultSplitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.SplitType == "" {
		writeError(w, "split_type required", http.StatusBadRequest)
		return
	}
	if len(req.SplitRatio) == 0 {
		writeError(w, "split_ratio required", http.StatusBadRequest)
		return
	}

	err := svc.SetDefaultSplitStrategy(groupID, req.SplitType, req.SplitRatio)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"message": "Default split strategy updated"}, http.StatusOK)
}

// ------------------------- Main -------------------------

func main() {
	store := NewStore()
	svc = NewService(store)

	// Example data setup
	log.Println("Setting up example data...")
	alice := svc.CreateUser("Alice", "alice@example.com")
	bob := svc.CreateUser("Bob", "bob@example.com")
	charlie := svc.CreateUser("Charlie", "charlie@example.com")

	tripGroup, _ := svc.CreateGroup("Trip to NYC", []string{alice.ID, bob.ID, charlie.ID}, alice.ID)

	// Example: Add expense with equal split
	expenseReq := &AddExpenseRequest{
		Title:        "Dinner at Restaurant",
		Description:  "Group dinner",
		Amount:       300.00,
		Payer:        alice.ID,
		CreatedBy:    alice.ID,
		GroupID:      &tripGroup.ID,
		ExpenseType:  ExpenseTypeEqual,
		SplitType:    SplitTypeEqual,
		Participants: []string{alice.ID, bob.ID, charlie.ID},
	}
	expense, _ := svc.AddExpense(expenseReq)
	log.Printf("Created expense: %s - $%.2f\n", expense.Title, expense.Amount)

	// Example: Add expense with percentage split
	expenseReq2 := &AddExpenseRequest{
		Title:       "Hotel Booking",
		Amount:      600.00,
		Payer:       bob.ID,
		CreatedBy:   bob.ID,
		GroupID:     &tripGroup.ID,
		ExpenseType: ExpenseTypePercentage,
		SplitType:   SplitTypePercentage,
		Percentages: map[string]float64{
			alice.ID:   50.0,
			bob.ID:     30.0,
			charlie.ID: 20.0,
		},
	}
	expense2, _ := svc.AddExpense(expenseReq2)
	log.Printf("Created expense: %s - $%.2f\n", expense2.Title, expense2.Amount)

	// Set up routes
	http.HandleFunc("/api/users/create", handleCreateUser)
	http.HandleFunc("/api/groups/create", handleCreateGroup)
	http.HandleFunc("/api/expenses", handleAddExpense)
	http.HandleFunc("/api/users/", handleGetBalances)
	http.HandleFunc("/api/groups/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/simplify-debts") {
			handleSimplifyDebts(w, r)
		} else if strings.HasSuffix(path, "/activity") {
			handleGetActivities(w, r)
		} else if strings.HasSuffix(path, "/default-split") {
			handleSetDefaultSplit(w, r)
		} else {
			writeError(w, "Not found", http.StatusNotFound)
		}
	})

	log.Println("Server running on :8080")
	log.Println("\nExample API endpoints:")
	log.Println("  POST   /api/users/create")
	log.Println("  POST   /api/groups/create")
	log.Println("  POST   /api/expenses")
	log.Println("  GET    /api/users/{user_id}/balances?group_id={group_id}")
	log.Println("  POST   /api/groups/{group_id}/simplify-debts")
	log.Println("  GET    /api/groups/{group_id}/activity")
	log.Println("  PUT    /api/groups/{group_id}/default-split")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
