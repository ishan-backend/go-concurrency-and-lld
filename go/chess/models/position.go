package models

// Position represents a square on the chess board
type Position struct {
	Row int // 0-7 (rank 1-8)
	Col int // 0-7 (file a-h)
}

// NewPosition creates a new position
func NewPosition(row, col int) Position {
	return Position{Row: row, Col: col}
}

// IsValid checks if the position is within board bounds
func (p Position) IsValid() bool {
	return p.Row >= 0 && p.Row < 8 && p.Col >= 0 && p.Col < 8
}

// ToAlgebraic converts position to algebraic notation (e.g., "e4")
func (p Position) ToAlgebraic() string {
	if !p.IsValid() {
		return ""
	}
	return string(rune('a'+p.Col)) + string(rune('1'+p.Row))
}

// FromAlgebraic creates a position from algebraic notation (e.g., "e4")
func FromAlgebraic(alg string) (Position, bool) {
	if len(alg) != 2 {
		return Position{}, false
	}
	col := int(alg[0] - 'a')
	row := int(alg[1] - '1')
	pos := NewPosition(row, col)
	if !pos.IsValid() {
		return Position{}, false
	}
	return pos, true
}

// Equals checks if two positions are equal
func (p Position) Equals(other Position) bool {
	return p.Row == other.Row && p.Col == other.Col
}

