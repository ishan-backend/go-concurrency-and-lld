package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// RookBehavior implements rook movement logic
type RookBehavior struct{}

func NewRookBehavior() *RookBehavior {
	return &RookBehavior{}
}

func (r *RookBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	// Directions: up, down, left, right
	directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	
	for _, dir := range directions {
		for i := 1; i < 8; i++ {
			newPos := models.NewPosition(pos.Row+i*dir[0], pos.Col+i*dir[1])
			if !newPos.IsValid() {
				break
			}
			
			if board.IsEmpty(newPos) {
				moves = append(moves, models.NewMove(pos, newPos))
			} else if board.IsOccupiedByColor(newPos, color.Opposite()) {
				moves = append(moves, models.NewMove(pos, newPos))
				break
			} else {
				break // Blocked by own piece
			}
		}
	}
	
	return moves
}

func (r *RookBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	moves := r.GetValidMoves(board, from, color)
	for _, move := range moves {
		if move.To.Equals(to) {
			return true
		}
	}
	return false
}

