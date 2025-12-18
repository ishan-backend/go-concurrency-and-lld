package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// BishopBehavior implements bishop movement logic
type BishopBehavior struct{}

func NewBishopBehavior() *BishopBehavior {
	return &BishopBehavior{}
}

func (b *BishopBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	// Directions: diagonals
	directions := [][]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	
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

func (b *BishopBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	moves := b.GetValidMoves(board, from, color)
	for _, move := range moves {
		if move.To.Equals(to) {
			return true
		}
	}
	return false
}

