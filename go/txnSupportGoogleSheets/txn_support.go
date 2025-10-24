package main

import (
	"errors"
	"fmt"
)

type Sheet struct {
	table map[string]map[string]string
	txn   *TransactionManager
}

type Change struct {
	row, col             string
	oldValue             *string
	newValue             *string
	rowAdded, rowDeleted bool
	colAdded, colDeleted bool
}

type TransactionManager struct {
	inTxn   bool
	changes []Change
}

func NewSheet() *Sheet {
	return &Sheet{
		table: make(map[string]map[string]string),
		txn:   &TransactionManager{},
	}
}

// ---------------------- Transaction Control ----------------------

func (s *Sheet) StartTxn() error {
	if s.txn.inTxn {
		return errors.New("transaction already in progress")
	}
	s.txn.inTxn = true
	s.txn.changes = nil
	return nil
}

func (s *Sheet) EndTxn() error {
	if !s.txn.inTxn {
		return errors.New("no transaction active")
	}
	s.txn.inTxn = false
	s.txn.changes = nil // discard undo history
	return nil
}

func (s *Sheet) Rollback() error {
	if !s.txn.inTxn {
		return errors.New("no transaction active")
	}

	// revert in reverse order
	for i := len(s.txn.changes) - 1; i >= 0; i-- {
		c := s.txn.changes[i]
		if c.rowAdded {
			delete(s.table, c.row)
		} else if c.rowDeleted {
			s.table[c.row] = make(map[string]string)
			if c.oldValue != nil && c.col != "" {
				s.table[c.row][c.col] = *c.oldValue
			}
		} else if c.colAdded {
			for _, row := range s.table {
				delete(row, c.col)
			}
		} else if c.colDeleted {
			for _, row := range s.table {
				if c.oldValue != nil {
					row[c.col] = *c.oldValue
				}
			}
		} else {
			if c.oldValue == nil {
				delete(s.table[c.row], c.col)
			} else {
				s.table[c.row][c.col] = *c.oldValue
			}
		}
	}

	s.txn.changes = nil
	s.txn.inTxn = false
	return nil
}

// ---------------------- Sheet Operations ----------------------

func (s *Sheet) AddRow(row string) {
	if _, ok := s.table[row]; !ok {
		s.table[row] = make(map[string]string)
		if s.txn.inTxn {
			s.txn.changes = append(s.txn.changes, Change{row: row, rowAdded: true})
		}
	}
}

func (s *Sheet) DeleteRow(row string) {
	if _, ok := s.table[row]; ok {
		if s.txn.inTxn {
			s.txn.changes = append(s.txn.changes, Change{row: row, rowDeleted: true})
		}
		delete(s.table, row)
	}
}

func (s *Sheet) AddCol(col string) {
	for row := range s.table {
		s.table[row][col] = ""
	}
	if s.txn.inTxn {
		s.txn.changes = append(s.txn.changes, Change{col: col, colAdded: true})
	}
}

func (s *Sheet) DeleteCol(col string) {
	for row := range s.table {
		val := s.table[row][col]
		if s.txn.inTxn {
			s.txn.changes = append(s.txn.changes, Change{row: row, col: col, oldValue: &val, colDeleted: true})
		}
		delete(s.table[row], col)
	}
}

func (s *Sheet) UpdateVal(row, col, val string) {
	oldVal, existed := s.table[row][col]
	if s.txn.inTxn {
		var oldPtr *string
		if existed {
			oldPtr = &oldVal
		}
		s.txn.changes = append(s.txn.changes, Change{
			row: row, col: col, oldValue: oldPtr, newValue: &val,
		})
	}
	if _, ok := s.table[row]; !ok {
		s.table[row] = make(map[string]string)
	}
	s.table[row][col] = val
}

func (s *Sheet) AddVal(row, col string, delta int) error {
	val, ok := s.table[row][col]
	if !ok {
		return errors.New("cell not found")
	}
	var curr int
	fmt.Sscanf(val, "%d", &curr)
	newVal := fmt.Sprintf("%d", curr+delta)
	s.UpdateVal(row, col, newVal)
	return nil
}

// ---------------------- Debug ----------------------

func (s *Sheet) Print() {
	fmt.Println("Sheet:")
	for r, cols := range s.table {
		fmt.Printf("%s -> %v\n", r, cols)
	}
}

func main() {
	s := NewSheet()
	s.AddRow("r1")
	s.AddCol("A")
	s.UpdateVal("r1", "A", "10")

	s.StartTxn()
	s.UpdateVal("r1", "A", "20")
	s.AddRow("r2")
	s.UpdateVal("r2", "A", "5")

	fmt.Println("Before rollback:")
	s.Print()

	s.Rollback()

	fmt.Println("\nAfter rollback:")
	s.Print()
}
