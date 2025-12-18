package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// PieceBehavior defines the interface for piece movement logic
type PieceBehavior interface {
	// GetValidMoves returns all valid moves for a piece at the given position
	GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move
	
	// CanMove checks if a piece can move from 'from' to 'to'
	CanMove(board *models.Board, from, to models.Position, color models.Color) bool
}

