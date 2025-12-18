package models

// Color represents the color of a chess piece
type Color int

const (
	White Color = iota
	Black
)

func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	default:
		return "Unknown"
	}
}

// Opposite returns the opposite color
func (c Color) Opposite() Color {
	if c == White {
		return Black
	}
	return White
}

