package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

// KingBehavior implements king movement logic
type KingBehavior struct{}

func NewKingBehavior() *KingBehavior {
	return &KingBehavior{}
}

func (k *KingBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
	var moves []models.Move
	
	// King moves one square in any direction
	offsets := [][]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
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
	
	// Castling (simplified - no check for check or moved pieces)
	if pos.Row == 0 || pos.Row == 7 {
		// Kingside castling
		if board.IsEmpty(models.NewPosition(pos.Row, 5)) &&
			board.IsEmpty(models.NewPosition(pos.Row, 6)) {
			rookPos := models.NewPosition(pos.Row, 7)
			if piece := board.GetPiece(rookPos); piece != nil &&
				piece.Type == models.Rook && piece.Color == color {
				moves = append(moves, models.NewMove(pos, models.NewPosition(pos.Row, 6)))
			}
		}
		
		// Queenside castling
		if board.IsEmpty(models.NewPosition(pos.Row, 1)) &&
			board.IsEmpty(models.NewPosition(pos.Row, 2)) &&
			board.IsEmpty(models.NewPosition(pos.Row, 3)) {
			rookPos := models.NewPosition(pos.Row, 0)
			if piece := board.GetPiece(rookPos); piece != nil &&
				piece.Type == models.Rook && piece.Color == color {
				moves = append(moves, models.NewMove(pos, models.NewPosition(pos.Row, 2)))
			}
		}
	}
	
	return moves
}

func (k *KingBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
	moves := k.GetValidMoves(board, from, color)
	for _, move := range moves {
		if move.To.Equals(to) {
			return true
		}
	}
	return false
}

