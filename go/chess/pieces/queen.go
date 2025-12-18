package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// QueenBehavior implements queen movement logic (combines rook and bishop)
type QueenBehavior struct {
	rook   *RookBehavior
	bishop *BishopBehavior
}

func NewQueenBehavior() *QueenBehavior {
	return &QueenBehavior{
		rook:   NewRookBehavior(),
		bishop: NewBishopBehavior(),
	}
}

func (q *QueenBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	// Queen moves like both rook and bishop
	moves = append(moves, q.rook.GetValidMoves(board, pos, color)...)
	moves = append(moves, q.bishop.GetValidMoves(board, pos, color)...)
	
	return moves
}

func (q *QueenBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	return q.rook.CanMove(board, from, to, color) || q.bishop.CanMove(board, from, to, color)
}

