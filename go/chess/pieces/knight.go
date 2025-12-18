package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// KnightBehavior implements knight movement logic
type KnightBehavior struct{}

func NewKnightBehavior() *KnightBehavior {
	return &KnightBehavior{}
}

func (k *KnightBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	// Knight moves in L-shape: 2 squares in one direction, 1 square perpendicular
	offsets := [][]int{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	
	for _, offset := range offsets {
		newPos := models.NewPosition(pos.Row+offset[0], pos.Col+offset[1])
		if !newPos.IsValid() {
			continue
		}
		
		if board.IsEmpty(newPos) || board.IsOccupiedByColor(newPos, color.Opposite()) {
			moves = append(moves, models.NewMove(pos, newPos))
		}
	}
	
	return moves
}

func (k *KnightBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	moves := k.GetValidMoves(board, from, color)
	for _, move := range moves {
		if move.To.Equals(to) {
			return true
		}
	}
	return false
}

