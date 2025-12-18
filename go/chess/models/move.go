package models

// Move represents a chess move
type Move struct {
	From      Position
	To        Position
	Promotion PieceType // For pawn promotion
}

// NewMove creates a new move
func NewMove(from, to Position) Move {
	return Move{
		From:      from,
		To:        to,
		Promotion: -1, // No promotion
	}
}

// NewPromotionMove creates a move with promotion
func NewPromotionMove(from, to Position, promotion PieceType) Move {
	return Move{
		From:      from,
		To:        to,
		Promotion: promotion,
	}
}

// ToAlgebraic converts move to algebraic notation (simplified)
func (m Move) ToAlgebraic() string {
	return m.From.ToAlgebraic() + m.To.ToAlgebraic()
}

