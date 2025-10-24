package main

import (
	"fmt"
	"sync"
)

// -------------------- Minimal (LeetCode-style) O(1) Implementation --------------------

type TicTacToe struct {
	n        int
	rows     []int
	cols     []int
	diag     int
	antiDiag int
}

// Constructor returns a ready TicTacToe for an n x n board.
func NewTicTacToe(n int) *TicTacToe {
	return &TicTacToe{
		n:    n,
		rows: make([]int, n),
		cols: make([]int, n),
	}
}

// Move places a mark by player (1 or 2) at (row, col).
// Returns 0 = no winner; 1 = player1 wins; 2 = player2 wins.
func (t *TicTacToe) Move(row, col, player int) int {
	var add int
	if player == 1 {
		add = 1
	} else {
		add = -1
	}

	t.rows[row] += add
	if abs(t.rows[row]) == t.n {
		if t.rows[row] > 0 {
			return 1
		}
		return 2
	}

	t.cols[col] += add
	if abs(t.cols[col]) == t.n {
		if t.cols[col] > 0 {
			return 1
		}
		return 2
	}

	if row == col {
		t.diag += add
		if abs(t.diag) == t.n {
			if t.diag > 0 {
				return 1
			}
			return 2
		}
	}

	if row+col == t.n-1 {
		t.antiDiag += add
		if abs(t.antiDiag) == t.n {
			if t.antiDiag > 0 {
				return 1
			}
			return 2
		}
	}

	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// -------------------- Safe version with board validation and optional concurrency --------------------

// TicTacToeSafe prevents placing in occupied cells. Returns:
// -1 => invalid move (occupied or out of bounds)
// 0 => no winner
// 1 or 2 => winner
type TicTacToeSafe struct {
	n        int
	rows     []int
	cols     []int
	diag     int
	antiDiag int
	board    [][]int    // 0 empty, 1 or 2
	mu       sync.Mutex // protect concurrent access (optional)
}

func NewTicTacToeSafe(n int) *TicTacToeSafe {
	board := make([][]int, n)
	for i := 0; i < n; i++ {
		board[i] = make([]int, n)
	}
	return &TicTacToeSafe{
		n:     n,
		rows:  make([]int, n),
		cols:  make([]int, n),
		board: board,
	}
}

// Thread-safe Move with validation
func (t *TicTacToeSafe) Move(row, col, player int) int {
	// bounds check
	if row < 0 || row >= t.n || col < 0 || col >= t.n || (player != 1 && player != 2) {
		return -1
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// occupied?
	if t.board[row][col] != 0 {
		return -1
	}

	t.board[row][col] = player

	var add int
	if player == 1 {
		add = 1
	} else {
		add = -1
	}

	t.rows[row] += add
	if abs(t.rows[row]) == t.n {
		if t.rows[row] > 0 {
			return 1
		}
		return 2
	}

	t.cols[col] += add
	if abs(t.cols[col]) == t.n {
		if t.cols[col] > 0 {
			return 1
		}
		return 2
	}

	if row == col {
		t.diag += add
		if abs(t.diag) == t.n {
			if t.diag > 0 {
				return 1
			}
			return 2
		}
	}

	if row+col == t.n-1 {
		t.antiDiag += add
		if abs(t.antiDiag) == t.n {
			if t.antiDiag > 0 {
				return 1
			}
			return 2
		}
	}

	return 0
}

// -------------------- Example usage --------------------

func main() {
	// Minimal example
	t := NewTicTacToe(3)
	fmt.Println(t.Move(0, 0, 1)) // 0
	fmt.Println(t.Move(0, 2, 2)) // 0
	fmt.Println(t.Move(2, 2, 1)) // 0
	fmt.Println(t.Move(1, 1, 2)) // 0
	fmt.Println(t.Move(2, 0, 1)) // 0
	fmt.Println(t.Move(1, 0, 2)) // 0
	fmt.Println(t.Move(2, 1, 1)) // 1 -> player1 wins

	// Safe example
	ts := NewTicTacToeSafe(3)
	fmt.Println(ts.Move(0, 0, 1)) // 0
	fmt.Println(ts.Move(0, 0, 2)) // -1 (occupied)
}
