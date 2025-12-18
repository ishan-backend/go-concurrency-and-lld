package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/chess/models"
	"github.com/ishanpandey/LowLevelDesign/go/chess/pieces"
)

// Game represents a chess game
type Game struct {
	board           *models.Board
	currentTurn     models.Color
	state           models.GameState
	moveHistory     []models.Move
	pieceBehaviors  map[models.PieceType]pieces.PieceBehavior
	mu              sync.RWMutex
}

// NewGame creates a new chess game
func NewGame() *Game {
	game := &Game{
		board:          models.NewStandardBoard(),
		currentTurn:    models.White,
		state:          models.Playing,
		moveHistory:    []models.Move{},
		pieceBehaviors: make(map[models.PieceType]pieces.PieceBehavior),
	}
	
	// Initialize piece behaviors
	game.pieceBehaviors[models.Pawn] = pieces.NewPawnBehavior()
	game.pieceBehaviors[models.Rook] = pieces.NewRookBehavior()
	game.pieceBehaviors[models.Knight] = pieces.NewKnightBehavior()
	game.pieceBehaviors[models.Bishop] = pieces.NewBishopBehavior()
	game.pieceBehaviors[models.Queen] = pieces.NewQueenBehavior()
	game.pieceBehaviors[models.King] = pieces.NewKingBehavior()
	
	return game
}

// GetBoard returns a copy of the current board
func (g *Game) GetBoard() *models.Board {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.board.Copy()
}

// GetCurrentTurn returns whose turn it is
func (g *Game) GetCurrentTurn() models.Color {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.currentTurn
}

// GetState returns the current game state
func (g *Game) GetState() models.GameState {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.state
}

// MakeMove attempts to make a move
func (g *Game) MakeMove(move models.Move) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.state != models.Playing {
		return errors.New("game is not in playing state")
	}
	
	// Validate move
	if err := g.validateMove(move); err != nil {
		return err
	}
	
	// Make the move
	capturedPiece := g.board.GetPiece(move.To)
	g.board.MakeMove(move)
	
	// Check if move puts own king in check
	if g.isInCheck(g.currentTurn) {
		// Undo move
		g.board.UnmakeMove(move, capturedPiece)
		return errors.New("move puts own king in check")
	}
	
	// Record move
	g.moveHistory = append(g.moveHistory, move)
	
	// Switch turn
	g.currentTurn = g.currentTurn.Opposite()
	
	// Check game state
	g.updateGameState()
	
	return nil
}

// validateMove validates if a move is legal
func (g *Game) validateMove(move models.Move) error {
	if !move.From.IsValid() || !move.To.IsValid() {
		return errors.New("invalid position")
	}
	
	piece := g.board.GetPiece(move.From)
	if piece == nil {
		return errors.New("no piece at source position")
	}
	
	if piece.Color != g.currentTurn {
		return errors.New("not your turn")
	}
	
	// Check if destination has own piece
	if g.board.IsOccupiedByColor(move.To, g.currentTurn) {
		return errors.New("cannot capture own piece")
	}
	
	// Get piece behavior and validate move
	behavior, exists := g.pieceBehaviors[piece.Type]
	if !exists {
		return errors.New("unknown piece type")
	}
	
	if !behavior.CanMove(g.board, move.From, move.To, piece.Color) {
		return errors.New("illegal move")
	}
	
	return nil
}

// isInCheck checks if the given color's king is in check
func (g *Game) isInCheck(color models.Color) bool {
	// Find the king
	var kingPos models.Position
	found := false
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := models.NewPosition(row, col)
			piece := g.board.GetPiece(pos)
			if piece != nil && piece.Type == models.King && piece.Color == color {
				kingPos = pos
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	
	if !found {
		return false // No king found (shouldn't happen in normal game)
	}
	
	// Check if any opponent piece can attack the king
	opponentColor := color.Opposite()
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := models.NewPosition(row, col)
			piece := g.board.GetPiece(pos)
			if piece != nil && piece.Color == opponentColor {
				behavior, exists := g.pieceBehaviors[piece.Type]
				if exists && behavior.CanMove(g.board, pos, kingPos, opponentColor) {
					return true
				}
			}
		}
	}
	
	return false
}

// hasValidMoves checks if the given color has any valid moves
func (g *Game) hasValidMoves(color models.Color) bool {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := models.NewPosition(row, col)
			piece := g.board.GetPiece(pos)
			if piece != nil && piece.Color == color {
				behavior, exists := g.pieceBehaviors[piece.Type]
				if exists {
					moves := behavior.GetValidMoves(g.board, pos, color)
					for _, move := range moves {
						// Try the move
						capturedPiece := g.board.GetPiece(move.To)
						g.board.MakeMove(move)
						
						// Check if it leaves king in check
						inCheck := g.isInCheck(color)
						
						// Undo move
						g.board.UnmakeMove(move, capturedPiece)
						
						if !inCheck {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// updateGameState updates the game state based on current position
func (g *Game) updateGameState() {
	// Check if current player is in checkmate or stalemate
	inCheck := g.isInCheck(g.currentTurn)
	hasMoves := g.hasValidMoves(g.currentTurn)
	
	if inCheck && !hasMoves {
		// Checkmate
		if g.currentTurn == models.White {
			g.state = models.BlackWon
		} else {
			g.state = models.WhiteWon
		}
	} else if !inCheck && !hasMoves {
		// Stalemate
		g.state = models.Stalemate
	}
	// Otherwise, game continues
}

// GetValidMoves returns all valid moves for a piece at the given position
func (g *Game) GetValidMoves(pos models.Position) []models.Move {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	piece := g.board.GetPiece(pos)
	if piece == nil || piece.Color != g.currentTurn {
		return []models.Move{}
	}
	
	behavior, exists := g.pieceBehaviors[piece.Type]
	if !exists {
		return []models.Move{}
	}
	
	allMoves := behavior.GetValidMoves(g.board, pos, piece.Color)
	var validMoves []models.Move
	
	for _, move := range allMoves {
		// Try the move
		capturedPiece := g.board.GetPiece(move.To)
		g.board.MakeMove(move)
		
		// Check if it leaves king in check
		inCheck := g.isInCheck(piece.Color)
		
		// Undo move
		g.board.UnmakeMove(move, capturedPiece)
		
		if !inCheck {
			validMoves = append(validMoves, move)
		}
	}
	
	return validMoves
}

// PrintBoard prints the current board state
func (g *Game) PrintBoard() {
	g.mu.RLock()
	defer g.mu.RUnlock()
	g.board.Print()
	fmt.Printf("\nTurn: %s | State: %s\n\n", g.currentTurn, g.state)
}

// IsInCheck checks if current player is in check
func (g *Game) IsInCheck() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.isInCheck(g.currentTurn)
}

