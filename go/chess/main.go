package main

import (
	"fmt"
	"log"

	"github.com/ishanpandey/LowLevelDesign/go/chess/models"
)

func main() {
	// Create a new game
	game := NewGame()
	
	fmt.Println("=== Chess Game ===")
	fmt.Println()
	
	// Print initial board
	game.PrintBoard()
	
	// Example moves
	moves := []string{
		"e2e4", // White pawn moves
		"e7e5", // Black pawn moves
		"g1f3", // White knight
		"b8c6", // Black knight
		"f1c4", // White bishop
		"f8c5", // Black bishop
		"e1g1", // White kingside castling (simplified)
	}
	
	for _, moveStr := range moves {
		from, ok1 := models.FromAlgebraic(moveStr[:2])
		to, ok2 := models.FromAlgebraic(moveStr[2:])
		
		if !ok1 || !ok2 {
			fmt.Printf("Invalid move format: %s\n", moveStr)
			continue
		}
		
		move := models.NewMove(from, to)
		
		fmt.Printf("Attempting move: %s\n", move.ToAlgebraic())
		
		if err := game.MakeMove(move); err != nil {
			log.Printf("Move failed: %v\n", err)
		} else {
			fmt.Printf("Move successful!\n")
			game.PrintBoard()
			
			if game.GetState() != models.Playing {
				fmt.Printf("Game Over: %s\n", game.GetState())
				break
			}
		}
		fmt.Println()
	}
	
	// Demonstrate getting valid moves
	fmt.Println("=== Valid Moves Example ===")
	pos, _ := models.FromAlgebraic("e2")
	validMoves := game.GetValidMoves(pos)
	fmt.Printf("Valid moves for e2: ")
	for _, move := range validMoves {
		fmt.Printf("%s ", move.ToAlgebraic())
	}
	fmt.Println()
}

