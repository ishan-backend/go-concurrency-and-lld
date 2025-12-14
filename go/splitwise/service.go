package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Service struct {
	store *Store
}

func NewService(st *Store) *Service {
	return &Service{store: st}
}

func (svc *Service) CreateUser(name, email string) *User {
	u := &User{
		ID:        newID("user"),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	svc.store.addUser(u)
	return u
}

func (svc *Service) CreateGroup(name string, members []string, admin string) (*Group, error) {
	// Validate admin is in members
	adminInMembers := false
	for _, m := range members {
		if m == admin {
			adminInMembers = true
			break
		}
	}
	if !adminInMembers {
		return nil, fmt.Errorf("admin must be a member of the group")
	}

	g := &Group{
		ID:                newID("grp"),
		Name:              name,
		Members:           members,
		Admin:             admin,
		DefaultSplitType:  SplitTypeEqual,
		DefaultSplitRatio: make(map[string]float64),
		CreatedAt:         time.Now(),
	}
	svc.store.addGroup(g)

	// Record activity
	svc.store.addActivity(&Activity{
		ID:          newID("act"),
		Type:        ActivityTypeGroupCreated,
		UserID:      admin,
		GroupID:     &g.ID,
		Description: fmt.Sprintf("%s created group '%s'", admin, name),
		CreatedAt:   time.Now(),
	})

	return g, nil
}

func (svc *Service) AddExpense(req *AddExpenseRequest) (*Expense, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Validate payer and created_by exist
	if _, ok := svc.store.getUser(req.Payer); !ok {
		return nil, fmt.Errorf("payer not found")
	}
	if _, ok := svc.store.getUser(req.CreatedBy); !ok {
		return nil, fmt.Errorf("created_by user not found")
	}

	// Validate group if provided
	if req.GroupID != nil {
		if _, ok := svc.store.getGroup(*req.GroupID); !ok {
			return nil, fmt.Errorf("group not found")
		}
	}

	// Calculate shares based on split type
	shares, err := svc.calculateShares(req)
	if err != nil {
		return nil, err
	}

	// Validate shares sum to amount
	sum := 0.0
	for _, v := range shares {
		sum += v
	}
	if math.Abs(sum-req.Amount) > 0.01 {
		return nil, fmt.Errorf("shares don't sum up to amount: got %.2f, expected %.2f", sum, req.Amount)
	}

	e := &Expense{
		ID:          newID("exp"),
		GroupID:     req.GroupID,
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		Payer:       req.Payer,
		CreatedBy:   req.CreatedBy,
		ExpenseType: req.ExpenseType,
		SplitType:   req.SplitType,
		Shares:      shares,
		CreatedAt:   time.Now(),
	}

	svc.store.recordExpense(e)

	// Record activity
	svc.store.addActivity(&Activity{
		ID:          newID("act"),
		Type:        ActivityTypeExpenseAdded,
		UserID:      req.CreatedBy,
		GroupID:     req.GroupID,
		ExpenseID:   &e.ID,
		Description: fmt.Sprintf("%s added expense '%s' of $%.2f", req.CreatedBy, req.Title, req.Amount),
		CreatedAt:   time.Now(),
	})

	return e, nil
}

func (svc *Service) calculateShares(req *AddExpenseRequest) (map[string]float64, error) {
	shares := make(map[string]float64)

	// If split type is not specified, try to infer from expense type
	splitType := req.SplitType
	if splitType == "" {
		splitType = SplitType(req.ExpenseType)
	}

	// If still empty and group has default, use group default
	if splitType == "" && req.GroupID != nil {
		if group, ok := svc.store.getGroup(*req.GroupID); ok {
			splitType = group.DefaultSplitType
			if len(req.Participants) == 0 {
				req.Participants = group.Members
			}
		}
	}

	// Default to EQUAL if still empty
	if splitType == "" {
		splitType = SplitTypeEqual
	}

	switch splitType {
	case SplitTypeEqual:
		if len(req.Participants) == 0 {
			return nil, fmt.Errorf("participants required for equal split")
		}
		amountPerPerson := req.Amount / float64(len(req.Participants))
		for _, p := range req.Participants {
			shares[p] = amountPerPerson
		}

	case SplitTypeUnequal:
		if len(req.ExactAmounts) == 0 {
			return nil, fmt.Errorf("exact_amounts required for unequal split")
		}
		shares = req.ExactAmounts

	case SplitTypePercentage:
		if len(req.Percentages) == 0 {
			// Try to use group default
			if req.GroupID != nil {
				if group, ok := svc.store.getGroup(*req.GroupID); ok && len(group.DefaultSplitRatio) > 0 {
					for userID, pct := range group.DefaultSplitRatio {
						shares[userID] = req.Amount * pct / 100.0
					}
				} else {
					return nil, fmt.Errorf("percentages required for percentage split")
				}
			} else {
				return nil, fmt.Errorf("percentages required for percentage split")
			}
		} else {
			// Validate percentages sum to 100
			sum := 0.0
			for _, pct := range req.Percentages {
				sum += pct
			}
			if math.Abs(sum-100.0) > 0.01 {
				return nil, fmt.Errorf("percentages must sum to 100, got %.2f", sum)
			}
			for userID, pct := range req.Percentages {
				shares[userID] = req.Amount * pct / 100.0
			}
		}

	case SplitTypeShares:
		if len(req.Shares) == 0 {
			// Try to use group default
			if req.GroupID != nil {
				if group, ok := svc.store.getGroup(*req.GroupID); ok && len(group.DefaultSplitRatio) > 0 {
					totalShares := 0.0
					for _, share := range group.DefaultSplitRatio {
						totalShares += share
					}
					for userID, share := range group.DefaultSplitRatio {
						shares[userID] = req.Amount * share / totalShares
					}
				} else {
					return nil, fmt.Errorf("shares required for share-based split")
				}
			} else {
				return nil, fmt.Errorf("shares required for share-based split")
			}
		} else {
			// Calculate total shares
			totalShares := 0.0
			for _, share := range req.Shares {
				totalShares += share
			}
			if totalShares <= 0 {
				return nil, fmt.Errorf("total shares must be positive")
			}
			for userID, share := range req.Shares {
				shares[userID] = req.Amount * share / totalShares
			}
		}

	default:
		return nil, fmt.Errorf("unknown split type: %s", splitType)
	}

	return shares, nil
}

func (svc *Service) GetBalances(userID string, groupID *string) (*BalanceResponse, error) {
	// Validate user exists
	if _, ok := svc.store.getUser(userID); !ok {
		return nil, fmt.Errorf("user not found")
	}

	// Validate group if provided
	if groupID != nil {
		if _, ok := svc.store.getGroup(*groupID); !ok {
			return nil, fmt.Errorf("group not found")
		}
	}

	netBalances := svc.store.getNetBalances(groupID)

	var balances []Balance
	if groupID != nil {
		balances = svc.store.getBalancesForUser(userID, groupID)
	} else {
		balances = svc.store.getBalancesForUser(userID, nil)
	}

	// Calculate totals
	totalOwed := 0.0
	totalOwedToMe := 0.0
	for _, balance := range balances {
		if balance.Amount > 0 {
			totalOwed += balance.Amount
		}
	}

	// Calculate what others owe to this user
	netBalance := netBalances[userID]
	if netBalance < 0 {
		totalOwedToMe = -netBalance
	}

	response := &BalanceResponse{
		UserID:        userID,
		GroupID:       groupID,
		TotalOwed:     totalOwed,
		TotalOwedToMe: totalOwedToMe,
		NetBalance:    netBalance,
	}

	if groupID != nil {
		response.GroupBalances = balances
	} else {
		response.GlobalBalances = balances
	}

	return response, nil
}

func (svc *Service) SimplifyDebts(groupID *string) (*SimplifyDebtsResponse, error) {
	netBalances := svc.store.getNetBalances(groupID)

	// Separate creditors and debtors
	type balanceEntry struct {
		uid string
		amt float64
	}

	var creditors []balanceEntry
	var debtors []balanceEntry

	for uid, balance := range netBalances {
		if math.Abs(balance) < 0.01 {
			continue
		}
		if balance > 0 {
			creditors = append(creditors, balanceEntry{uid: uid, amt: balance})
		} else {
			debtors = append(debtors, balanceEntry{uid: uid, amt: -balance})
		}
	}

	// Sort by amount (largest first)
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].amt > creditors[j].amt
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].amt > debtors[j].amt
	})

	// Greedy matching
	var transactions []Transaction
	ci, di := 0, 0

	for ci < len(creditors) && di < len(debtors) {
		c := creditors[ci]
		d := debtors[di]

		settleAmount := math.Min(c.amt, d.amt)

		tr, err := svc.store.settle(d.uid, c.uid, settleAmount, groupID)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, *tr)

		// Record activity
		svc.store.addActivity(&Activity{
			ID:          newID("act"),
			Type:        ActivityTypeSettlement,
			UserID:      d.uid,
			GroupID:     groupID,
			Description: fmt.Sprintf("%s settled $%.2f to %s", d.uid, settleAmount, c.uid),
			CreatedAt:   time.Now(),
		})

		// Update amounts
		c.amt -= settleAmount
		d.amt -= settleAmount

		if c.amt <= 0.01 {
			ci++
		}
		if d.amt <= 0.01 {
			di++
		}
	}

	// Count original transactions (approximate - number of non-zero balances)
	originalCount := len(creditors) + len(debtors)

	return &SimplifyDebtsResponse{
		SimplifiedTransactions:     transactions,
		OriginalTransactionCount:   originalCount,
		SimplifiedTransactionCount: len(transactions),
	}, nil
}

func (svc *Service) GetActivities(groupID *string, limit, offset int) (*ActivityResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	activities := svc.store.getActivities(groupID, limit, offset)

	// Count total
	total := len(svc.store.getActivities(groupID, 10000, 0))

	return &ActivityResponse{
		Activities: convertActivities(activities),
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	}, nil
}

func convertActivities(activities []*Activity) []Activity {
	result := make([]Activity, len(activities))
	for i, a := range activities {
		result[i] = *a
	}
	return result
}

func (svc *Service) SetDefaultSplitStrategy(groupID string, splitType SplitType, splitRatio map[string]float64) error {
	// Validate group exists
	group, ok := svc.store.getGroup(groupID)
	if !ok {
		return fmt.Errorf("group not found")
	}

	// Validate split ratio sums correctly
	if splitType == SplitTypePercentage {
		sum := 0.0
		for _, pct := range splitRatio {
			sum += pct
		}
		if math.Abs(sum-100.0) > 0.01 {
			return fmt.Errorf("percentages must sum to 100, got %.2f", sum)
		}
	}

	// Validate all users in split ratio are group members
	for userID := range splitRatio {
		isMember := false
		for _, member := range group.Members {
			if member == userID {
				isMember = true
				break
			}
		}
		if !isMember {
			return fmt.Errorf("user %s is not a member of the group", userID)
		}
	}

	return svc.store.updateGroupDefaultSplit(groupID, splitType, splitRatio)
}
