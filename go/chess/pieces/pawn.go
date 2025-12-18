package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// PawnBehavior implements pawn movement logic
type PawnBehavior struct{}

func NewPawnBehavior() *PawnBehavior {
	return &PawnBehavior{}
}

func (p *PawnBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	direction := 1 // White moves up (increasing row)
	if color == models.Black {
		direction = -1 // Black moves down (decreasing row)
	}
	
	startRow := 1 // White pawns start at row 1 (rank 2)
	if color == models.Black {
		startRow = 6 // Black pawns start at row 6 (rank 7)
	}
	
	// Move forward one square
	oneForward := models.NewPosition(pos.Row+direction, pos.Col)
	if oneForward.IsValid() && board.IsEmpty(oneForward) {
		moves = append(moves, models.NewMove(pos, oneForward))
		
		// Move forward two squares from starting position
		if pos.Row == startRow {
			twoForward := models.NewPosition(pos.Row+2*direction, pos.Col)
			if twoForward.IsValid() && board.IsEmpty(twoForward) {
				moves = append(moves, models.NewMove(pos, twoForward))
			}
		}
	}
	
	// Capture diagonally
	for _, colOffset := range []int{-1, 1} {
		capturePos := models.NewPosition(pos.Row+direction, pos.Col+colOffset)
		if capturePos.IsValid() {
			piece := board.GetPiece(capturePos)
			if piece != nil && piece.Color == color.Opposite() {
				moves = append(moves, models.NewMove(pos, capturePos))
			}
		}
	}
	
	// Promotion moves (simplified - always promote to queen)
	// White promotes at row 7 (rank 8), Black promotes at row 0 (rank 1)
	if (color == models.White && pos.Row == 6) || (color == models.Black && pos.Row == 1) {
		for i := range moves {
			if (color == models.White && moves[i].To.Row == 7) || (color == models.Black && moves[i].To.Row == 0) {
				moves[i].Promotion = models.Queen
			}
		}
	}
	
	return moves
}

func (p *PawnBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	moves := p.GetValidMoves(board, from, color)
	for _, move := range moves {
		if move.To.Equals(to) {
			return true
		}
	}
	return false
}

