package main

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu           sync.RWMutex
	users        map[string]*User
	groups       map[string]*Group
	expenses     map[string]*Expense
	transactions []*Transaction
	balances     map[string]map[string]*Balance // user_id -> (owed_to_user_id -> Balance)
	activities   []*Activity
}

func NewStore() *Store {
	return &Store{
		users:      make(map[string]*User),
		groups:     make(map[string]*Group),
		expenses:   make(map[string]*Expense),
		balances:   make(map[string]map[string]*Balance),
		activities: make([]*Activity, 0),
	}
}

func (s *Store) addUser(u *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
	if _, ok := s.balances[u.ID]; !ok {
		s.balances[u.ID] = make(map[string]*Balance)
	}
}

func (s *Store) getUser(userID string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	return user, ok
}

func (s *Store) addGroup(g *Group) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups[g.ID] = g
	for _, m := range g.Members {
		if _, ok := s.balances[m]; !ok {
			s.balances[m] = make(map[string]*Balance)
		}
		// Initialize balances between all members
		for _, other := range g.Members {
			if m != other {
				if _, ok := s.balances[m][other]; !ok {
					s.balances[m][other] = &Balance{
						UserID:  m,
						OwedTo:  other,
						Amount:  0,
						GroupID: &g.ID,
					}
				}
			}
		}
	}
}

func (s *Store) getGroup(groupID string) (*Group, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	group, ok := s.groups[groupID]
	return group, ok
}

func (s *Store) recordExpense(e *Expense) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expenses[e.ID] = e

	// Update balances
	for uid, shareAmount := range e.Shares {
		if uid == e.Payer {
			continue // Payer doesn't owe themselves
		}

		// Initialize balance if it doesn't exist
		if _, ok := s.balances[uid]; !ok {
			s.balances[uid] = make(map[string]*Balance)
		}
		if _, ok := s.balances[uid][e.Payer]; !ok {
			s.balances[uid][e.Payer] = &Balance{
				UserID:  uid,
				OwedTo:  e.Payer,
				Amount:  0,
				GroupID: e.GroupID,
			}
		}

		// Update balance: uid owes e.Payer shareAmount
		s.balances[uid][e.Payer].Amount += shareAmount

		// Also update reverse balance (payer is owed by uid)
		if _, ok := s.balances[e.Payer]; !ok {
			s.balances[e.Payer] = make(map[string]*Balance)
		}
		if _, ok := s.balances[e.Payer][uid]; !ok {
			s.balances[e.Payer][uid] = &Balance{
				UserID:  e.Payer,
				OwedTo:  uid,
				Amount:  0,
				GroupID: e.GroupID,
			}
		}
		s.balances[e.Payer][uid].Amount -= shareAmount
	}
}

func (s *Store) settle(from, to string, amount float64, groupID *string) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize balances if they don't exist
	if _, ok := s.balances[from]; !ok {
		s.balances[from] = make(map[string]*Balance)
	}
	if _, ok := s.balances[from][to]; !ok {
		s.balances[from][to] = &Balance{
			UserID:  from,
			OwedTo:  to,
			Amount:  0,
			GroupID: groupID,
		}
	}

	if _, ok := s.balances[to]; !ok {
		s.balances[to] = make(map[string]*Balance)
	}
	if _, ok := s.balances[to][from]; !ok {
		s.balances[to][from] = &Balance{
			UserID:  to,
			OwedTo:  from,
			Amount:  0,
			GroupID: groupID,
		}
	}

	// Update balances
	s.balances[from][to].Amount -= amount
	s.balances[to][from].Amount += amount

	tr := &Transaction{
		ID:        newID("txn"),
		From:      from,
		To:        to,
		Amount:    amount,
		GroupID:   groupID,
		CreatedAt: time.Now(),
	}
	s.transactions = append(s.transactions, tr)
	return tr, nil
}

func (s *Store) getNetBalances(groupID *string) map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make(map[string]float64)
	for from, balances := range s.balances {
		for to, balance := range balances {
			// Filter by group if specified
			if groupID != nil {
				// For group-specific, only include balances for this group
				if balance.GroupID == nil || *balance.GroupID != *groupID {
					continue
				}
			}
			// For global (groupID == nil), include all balances

			res[from] -= balance.Amount
			res[to] += balance.Amount
		}
	}
	return res
}

func (s *Store) getBalancesForUser(userID string, groupID *string) []Balance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Balance
	if userBalances, ok := s.balances[userID]; ok {
		for _, balance := range userBalances {
			// Filter by group if specified
			if groupID != nil {
				if balance.GroupID == nil || *balance.GroupID != *groupID {
					continue
				}
			}
			// For global (groupID == nil), include all balances
			// Only include positive amounts (what user owes)
			if balance.Amount > 0.01 {
				result = append(result, *balance)
			}
		}
	}
	return result
}

func (s *Store) addActivity(activity *Activity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activities = append(s.activities, activity)
}

func (s *Store) getActivities(groupID *string, limit, offset int) []*Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*Activity
	for _, act := range s.activities {
		if groupID != nil {
			if act.GroupID == nil || *act.GroupID != *groupID {
				continue
			}
		}
		filtered = append(filtered, act)
	}

	// Sort by CreatedAt descending (newest first)
	// Simple reverse iteration
	start := len(filtered) - offset - 1
	end := start - limit + 1
	if end < 0 {
		end = 0
	}
	if start < 0 {
		return []*Activity{}
	}

	var result []*Activity
	for i := start; i >= end && i >= 0; i-- {
		result = append(result, filtered[i])
	}

	return result
}

func (s *Store) updateGroupDefaultSplit(groupID string, splitType SplitType, splitRatio map[string]float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}

	group.DefaultSplitType = splitType
	group.DefaultSplitRatio = splitRatio
	return nil
}
