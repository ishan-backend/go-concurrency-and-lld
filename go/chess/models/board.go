package models

import "fmt"

// Board represents the chess board
type Board struct {
	squares [8][8]*Piece
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	return &Board{
		squares: [8][8]*Piece{},
	}
}

// NewStandardBoard creates a board with pieces in starting position
func NewStandardBoard() *Board {
	board := NewBoard()
	
	// Place pawns
	// White pawns at row 1 (rank 2), Black pawns at row 6 (rank 7)
	for col := 0; col < 8; col++ {
		board.SetPiece(NewPosition(1, col), NewPiece(Pawn, White))
		board.SetPiece(NewPosition(6, col), NewPiece(Pawn, Black))
	}
	
	// Place back row pieces
	// White back row at row 0 (rank 1), Black back row at row 7 (rank 8)
	backRow := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for col := 0; col < 8; col++ {
		board.SetPiece(NewPosition(0, col), NewPiece(backRow[col], White))
		board.SetPiece(NewPosition(7, col), NewPiece(backRow[col], Black))
	}
	
	return board
}

// GetPiece returns the piece at the given position
func (b *Board) GetPiece(pos Position) *Piece {
	if !pos.IsValid() {
		return nil
	}
	return b.squares[pos.Row][pos.Col]
}

// SetPiece places a piece at the given position
func (b *Board) SetPiece(pos Position, piece *Piece) {
	if pos.IsValid() {
		b.squares[pos.Row][pos.Col] = piece
	}
}

// RemovePiece removes a piece from the given position
func (b *Board) RemovePiece(pos Position) {
	if pos.IsValid() {
		b.squares[pos.Row][pos.Col] = nil
	}
}

// IsEmpty checks if a square is empty
func (b *Board) IsEmpty(pos Position) bool {
	return b.GetPiece(pos) == nil
}

// IsOccupiedByColor checks if a square is occupied by a piece of the given color
func (b *Board) IsOccupiedByColor(pos Position, color Color) bool {
	piece := b.GetPiece(pos)
	return piece != nil && piece.Color == color
}

// Copy creates a deep copy of the board
func (b *Board) Copy() *Board {
	newBoard := NewBoard()
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := b.squares[row][col]
			if piece != nil {
				newBoard.squares[row][col] = NewPiece(piece.Type, piece.Color)
			}
		}
	}
	return newBoard
}

// MakeMove applies a move to the board (without validation)
func (b *Board) MakeMove(move Move) {
	piece := b.GetPiece(move.From)
	if piece == nil {
		return
	}
	
	// Handle promotion
	if move.Promotion >= 0 {
		piece = NewPiece(move.Promotion, piece.Color)
	}
	
	b.RemovePiece(move.From)
	b.SetPiece(move.To, piece)
}

// UnmakeMove reverses a move (for move generation)
func (b *Board) UnmakeMove(move Move, capturedPiece *Piece) {
	piece := b.GetPiece(move.To)
	if piece == nil {
		return
	}
	
	// If it was a promotion, restore pawn
	if move.Promotion >= 0 {
		piece = NewPiece(Pawn, piece.Color)
	}
	
	b.RemovePiece(move.To)
	b.SetPiece(move.From, piece)
	if capturedPiece != nil {
		b.SetPiece(move.To, capturedPiece)
	}
}

// Print prints the board to console
func (b *Board) Print() {
	fmt.Println("  a b c d e f g h")
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d ", row+1)
		for col := 0; col < 8; col++ {
			piece := b.squares[row][col]
			if piece == nil {
				if (row+col)%2 == 0 {
					fmt.Print("· ")
				} else {
					fmt.Print("  ")
				}
			} else {
				fmt.Printf("%s ", piece.GetSymbol())
			}
		}
		fmt.Printf("%d\n", row+1)
	}
	fmt.Println("  a b c d e f g h")
}

