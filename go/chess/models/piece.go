package models

// PieceType represents the type of chess piece
type PieceType int

const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	Queen
	King
)

func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "Pawn"
	case Rook:
		return "Rook"
	case Knight:
		return "Knight"
	case Bishop:
		return "Bishop"
	case Queen:
		return "Queen"
	case King:
		return "King"
	default:
		return "Unknown"
	}
}

// Piece represents a chess piece
type Piece struct {
	Type  PieceType
	Color Color
}

// NewPiece creates a new piece
func NewPiece(pieceType PieceType, color Color) *Piece {
	return &Piece{
		Type:  pieceType,
		Color: color,
	}
}

// GetSymbol returns the Unicode symbol for the piece
func (p *Piece) GetSymbol() string {
	symbols := map[Color]map[PieceType]string{
		White: {
			Pawn:   "♙",
			Rook:   "♖",
			Knight: "♘",
			Bishop: "♗",
			Queen:  "♕",
			King:   "♔",
		},
		Black: {
			Pawn:   "♟",
			Rook:   "♜",
			Knight: "♞",
			Bishop: "♝",
			Queen:  "♛",
			King:   "♚",
		},
	}
	return symbols[p.Color][p.Type]
}

// GetFENChar returns the FEN character for the piece
func (p *Piece) GetFENChar() string {
	chars := map[Color]map[PieceType]string{
		White: {
			Pawn:   "P",
			Rook:   "R",
			Knight: "N",
			Bishop: "B",
			Queen:  "Q",
			King:   "K",
		},
		Black: {
			Pawn:   "p",
			Rook:   "r",
			Knight: "n",
			Bishop: "b",
			Queen:  "q",
			King:   "k",
		},
	}
	return chars[p.Color][p.Type]
}

