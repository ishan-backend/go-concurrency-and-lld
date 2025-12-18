package models

// GameState represents the current state of the game
type GameState int

const (
	Playing GameState = iota
	WhiteWon
	BlackWon
	Stalemate
	Draw
)

func (gs GameState) String() string {
	switch gs {
	case Playing:
		return "Playing"
	case WhiteWon:
		return "White Won"
	case BlackWon:
		return "Black Won"
	case Stalemate:
		return "Stalemate"
	case Draw:
		return "Draw"
	default:
		return "Unknown"
	}
}

