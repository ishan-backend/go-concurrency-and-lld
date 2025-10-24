package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ========== Core Entities ==========

// Position represents a position on the board
type Position int

// Snake represents a snake on the board
type Snake struct {
	Head Position
	Tail Position
}

// NewSnake creates a new snake
func NewSnake(head, tail Position) *Snake {
	if head <= tail {
		panic("Snake head must be greater than tail")
	}
	return &Snake{
		Head: head,
		Tail: tail,
	}
}

// Ladder represents a ladder on the board
type Ladder struct {
	Start Position
	End   Position
}

// NewLadder creates a new ladder
func NewLadder(start, end Position) *Ladder {
	if start >= end {
		panic("Ladder start must be less than end")
	}
	return &Ladder{
		Start: start,
		End:   end,
	}
}

// Player represents a player in the game
type Player struct {
	ID       string
	Name     string
	Position Position
}

// NewPlayer creates a new player
func NewPlayer(id, name string) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Position: 0, // Start position
	}
}

// Dice represents a dice
type Dice struct {
	Faces int
	rand  *rand.Rand
}

// NewDice creates a new dice
func NewDice(faces int) *Dice {
	return &Dice{
		Faces: faces,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Roll returns a random number between 1 and number of faces
func (d *Dice) Roll() int {
	return d.rand.Intn(d.Faces) + 1
}

// ========== Board ==========

// Board represents the game board
type Board struct {
	Size    int
	Snakes  map[Position]*Snake
	Ladders map[Position]*Ladder
}

// NewBoard creates a new board
func NewBoard(size int) *Board {
	return &Board{
		Size:    size,
		Snakes:  make(map[Position]*Snake),
		Ladders: make(map[Position]*Ladder),
	}
}

// AddSnake adds a snake to the board
func (b *Board) AddSnake(snake *Snake) error {
	if snake.Head > Position(b.Size) || snake.Tail < 1 {
		return fmt.Errorf("invalid snake position")
	}
	if _, exists := b.Snakes[snake.Head]; exists {
		return fmt.Errorf("snake already exists at position %d", snake.Head)
	}
	if _, exists := b.Ladders[snake.Head]; exists {
		return fmt.Errorf("ladder already exists at position %d", snake.Head)
	}
	b.Snakes[snake.Head] = snake
	return nil
}

// AddLadder adds a ladder to the board
func (b *Board) AddLadder(ladder *Ladder) error {
	if ladder.End > Position(b.Size) || ladder.Start < 1 {
		return fmt.Errorf("invalid ladder position")
	}
	if _, exists := b.Ladders[ladder.Start]; exists {
		return fmt.Errorf("ladder already exists at position %d", ladder.Start)
	}
	if _, exists := b.Snakes[ladder.Start]; exists {
		return fmt.Errorf("snake already exists at position %d", ladder.Start)
	}
	b.Ladders[ladder.Start] = ladder
	return nil
}

// GetNewPosition returns the new position after considering snakes and ladders
func (b *Board) GetNewPosition(position Position) Position {
	// Check for ladder
	if ladder, exists := b.Ladders[position]; exists {
		fmt.Printf("    🪜 Ladder! Climbing from %d to %d\n", position, ladder.End)
		return ladder.End
	}

	// Check for snake
	if snake, exists := b.Snakes[position]; exists {
		fmt.Printf("    🐍 Snake! Sliding from %d to %d\n", position, snake.Tail)
		return snake.Tail
	}

	return position
}

// Display prints the board configuration
func (b *Board) Display() {
	fmt.Println("\n📋 Board Configuration:")
	fmt.Printf("Size: %d\n", b.Size)

	fmt.Println("\n🐍 Snakes:")
	for head, snake := range b.Snakes {
		fmt.Printf("  %d → %d\n", head, snake.Tail)
	}

	fmt.Println("\n🪜 Ladders:")
	for start, ladder := range b.Ladders {
		fmt.Printf("  %d → %d\n", start, ladder.End)
	}
}

// ========== Game ==========

// GameStatus represents the current status of the game
type GameStatus int

const (
	NotStarted GameStatus = iota
	InProgress
	Finished
)

// Game represents the Snake and Ladder game
type Game struct {
	Board         *Board
	Players       []*Player
	Dice          *Dice
	CurrentPlayer int
	Status        GameStatus
	Winner        *Player
}

// NewGame creates a new game
func NewGame(boardSize int, dice *Dice) *Game {
	return &Game{
		Board:         NewBoard(boardSize),
		Players:       make([]*Player, 0),
		Dice:          dice,
		CurrentPlayer: 0,
		Status:        NotStarted,
		Winner:        nil,
	}
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(player *Player) error {
	if g.Status != NotStarted {
		return fmt.Errorf("cannot add player after game has started")
	}
	g.Players = append(g.Players, player)
	return nil
}

// Start starts the game
func (g *Game) Start() error {
	if len(g.Players) < 2 {
		return fmt.Errorf("need at least 2 players to start")
	}
	if g.Status != NotStarted {
		return fmt.Errorf("game already started")
	}
	g.Status = InProgress
	fmt.Println("\n🎮 Game Started!")
	return nil
}

// PlayTurn plays a turn for the current player
func (g *Game) PlayTurn() error {
	if g.Status != InProgress {
		return fmt.Errorf("game is not in progress")
	}

	player := g.Players[g.CurrentPlayer]
	diceValue := g.Dice.Roll()

	fmt.Printf("\n👤 %s's turn (Current Position: %d)\n", player.Name, player.Position)
	fmt.Printf("🎲 Rolled: %d\n", diceValue)

	newPosition := player.Position + Position(diceValue)

	// Check if new position exceeds board size
	if newPosition > Position(g.Board.Size) {
		fmt.Printf("    ❌ Cannot move! Position %d exceeds board size %d\n",
			newPosition, g.Board.Size)
		g.nextPlayer()
		return nil
	}

	// Move player
	player.Position = newPosition
	fmt.Printf("    Moved to: %d\n", player.Position)

	// Check for snake or ladder
	finalPosition := g.Board.GetNewPosition(player.Position)
	player.Position = finalPosition

	// Check for winner
	if player.Position == Position(g.Board.Size) {
		g.Status = Finished
		g.Winner = player
		fmt.Printf("\n🎉 %s WINS! 🎉\n", player.Name)
		return nil
	}

	fmt.Printf("    Final Position: %d\n", player.Position)

	g.nextPlayer()
	return nil
}

// nextPlayer moves to the next player
func (g *Game) nextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
}

// GetCurrentPlayer returns the current player
func (g *Game) GetCurrentPlayer() *Player {
	if len(g.Players) == 0 {
		return nil
	}
	return g.Players[g.CurrentPlayer]
}

// GetWinner returns the winner
func (g *Game) GetWinner() *Player {
	return g.Winner
}

// IsGameOver checks if the game is over
func (g *Game) IsGameOver() bool {
	return g.Status == Finished
}

// DisplayPlayerPositions displays all player positions
func (g *Game) DisplayPlayerPositions() {
	fmt.Println("\n📍 Player Positions:")
	for _, player := range g.Players {
		fmt.Printf("  %s: %d\n", player.Name, player.Position)
	}
}

// ========== Game Builder ==========

// GameBuilder helps in building a game with fluent interface
type GameBuilder struct {
	game *Game
}

// NewGameBuilder creates a new game builder
func NewGameBuilder(boardSize int, diceFaces int) *GameBuilder {
	return &GameBuilder{
		game: NewGame(boardSize, NewDice(diceFaces)),
	}
}

// WithSnake adds a snake to the game
func (gb *GameBuilder) WithSnake(head, tail Position) *GameBuilder {
	snake := NewSnake(head, tail)
	gb.game.Board.AddSnake(snake)
	return gb
}

// WithLadder adds a ladder to the game
func (gb *GameBuilder) WithLadder(start, end Position) *GameBuilder {
	ladder := NewLadder(start, end)
	gb.game.Board.AddLadder(ladder)
	return gb
}

// WithPlayer adds a player to the game
func (gb *GameBuilder) WithPlayer(id, name string) *GameBuilder {
	player := NewPlayer(id, name)
	gb.game.AddPlayer(player)
	return gb
}

// Build returns the built game
func (gb *GameBuilder) Build() *Game {
	return gb.game
}

// ========== Main ==========

func main() {
	fmt.Println("🐍🪜 SNAKE AND LADDER GAME 🪜🐍")
	fmt.Println("================================")

	// Build game using builder pattern
	game := NewGameBuilder(100, 6).
		// Add Snakes
		WithSnake(99, 54).
		WithSnake(70, 55).
		WithSnake(52, 42).
		WithSnake(25, 2).
		WithSnake(95, 72).
		// Add Ladders
		WithLadder(6, 25).
		WithLadder(11, 40).
		WithLadder(60, 85).
		WithLadder(46, 90).
		WithLadder(17, 69).
		// Add Players
		WithPlayer("p1", "Alice").
		WithPlayer("p2", "Bob").
		WithPlayer("p3", "Charlie").
		Build()

	// Display board
	game.Board.Display()

	// Start game
	err := game.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Play game
	turnCount := 0
	maxTurns := 100 // Prevent infinite games

	for !game.IsGameOver() && turnCount < maxTurns {
		err := game.PlayTurn()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		turnCount++

		// Display positions every 5 turns
		if turnCount%5 == 0 {
			game.DisplayPlayerPositions()
		}
	}

	// Display final results
	fmt.Println("\n" + "================================")
	fmt.Println("🏁 GAME OVER 🏁")
	fmt.Println("================================")

	if game.GetWinner() != nil {
		fmt.Printf("🏆 Winner: %s\n", game.GetWinner().Name)
	} else {
		fmt.Println("Game ended without a winner")
	}

	fmt.Printf("Total Turns: %d\n", turnCount)
	game.DisplayPlayerPositions()
}
