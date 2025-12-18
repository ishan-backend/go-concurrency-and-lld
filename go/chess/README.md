# Chess Game - Low Level Design

A comprehensive, extensible chess game implementation in Go that demonstrates clean architecture, design patterns, and proper separation of concerns.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Project Structure](#project-structure)
4. [Core Components](#core-components)
5. [Piece Movement Logic](#piece-movement-logic)
6. [Game Rules Implementation](#game-rules-implementation)
7. [Usage Examples](#usage-examples)
8. [Interview Explanation Guide](#interview-explanation-guide)
9. [Design Decisions](#design-decisions)
10. [Extensibility](#extensibility)

---

## Overview

This chess implementation provides a complete chess game engine with:
- **All 6 piece types** (Pawn, Rook, Knight, Bishop, Queen, King)
- **Move validation** with check/checkmate detection
- **Special moves** (Castling, Pawn Promotion)
- **Thread-safe** operations using mutexes
- **Clean architecture** with separation of concerns
- **Extensible design** using Strategy Pattern for piece behaviors

### Key Features

✅ **Strategy Pattern**: Each piece type has its own behavior implementation  
✅ **Thread-Safe**: All game operations are concurrent-safe using mutexes  
✅ **Extensible**: Easy to add new piece types or modify movement rules  
✅ **Clean Architecture**: Separation of models, pieces, and game logic  
✅ **Complete Rules**: Check, checkmate, stalemate detection  
✅ **Move Validation**: Prevents illegal moves and moves that put own king in check  

---

## Architecture & Design Patterns

### 1. **Strategy Pattern** ⭐ (Core Pattern)

**Purpose**: Encapsulate piece movement algorithms and make them interchangeable.

**Implementation**:
```go
type PieceBehavior interface {
    GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move
    CanMove(board *models.Board, from, to models.Position, color models.Color) bool
}
```

**Why it's powerful**:
- ✅ Each piece type has its own movement logic
- ✅ Easy to modify rules for specific pieces
- ✅ Open/Closed Principle: Add new pieces without modifying existing code
- ✅ Testable: Each behavior can be tested independently

**Real-world analogy**: Think of it like different vehicle types - a car, bicycle, and airplane all have different movement behaviors, but they share a common interface.

### 2. **Model-View Separation**

- **Models**: Core domain entities (Board, Piece, Position, Move)
- **Pieces**: Behavior implementations (Strategy Pattern)
- **Game**: Orchestration and game logic

### 3. **Immutable Operations**

The board provides `Copy()` method for creating snapshots, enabling safe move generation and validation without side effects.

---

## Project Structure

```
go/chess/
├── models/
│   ├── color.go           # Color enum (White/Black)
│   ├── position.go        # Board position with algebraic notation
│   ├── piece.go           # Piece type and representation
│   ├── move.go            # Move representation
│   ├── game_state.go      # Game state enum
│   └── board.go           # Board representation and operations
├── pieces/
│   ├── piece_interface.go # PieceBehavior interface
│   ├── pawn.go            # Pawn movement logic
│   ├── rook.go            # Rook movement logic
│   ├── knight.go          # Knight movement logic
│   ├── bishop.go          # Bishop movement logic
│   ├── queen.go           # Queen movement logic
│   └── king.go            # King movement logic (includes castling)
├── game.go                # Main game engine
├── main.go                # Demo and examples
├── go.mod                 # Go module definition
└── README.md              # This file
```

### Package Responsibilities

- **`models`**: Domain entities and data structures
- **`pieces`**: Piece movement behavior implementations (Strategy Pattern)
- **Root package**: Game engine, orchestration, and demo

---

## Core Components

### 1. Position (`models/position.go`)

Represents a square on the chess board.

```go
type Position struct {
    Row int // 0-7 (rank 1-8)
    Col int // 0-7 (file a-h)
}
```

**Key Methods**:
- `IsValid()`: Checks if position is within board bounds
- `ToAlgebraic()`: Converts to algebraic notation (e.g., "e4")
- `FromAlgebraic()`: Creates position from algebraic notation

### 2. Piece (`models/piece.go`)

Represents a chess piece with type and color.

```go
type Piece struct {
    Type  PieceType  // Pawn, Rook, Knight, Bishop, Queen, King
    Color Color      // White or Black
}
```

**Key Methods**:
- `GetSymbol()`: Returns Unicode chess symbol
- `GetFENChar()`: Returns FEN notation character

### 3. Board (`models/board.go`)

Represents the 8x8 chess board.

**Key Methods**:
- `GetPiece(pos)`: Returns piece at position
- `SetPiece(pos, piece)`: Places piece on board
- `IsEmpty(pos)`: Checks if square is empty
- `IsOccupiedByColor(pos, color)`: Checks if square has piece of color
- `Copy()`: Creates deep copy of board
- `MakeMove(move)`: Applies move to board
- `UnmakeMove(move, capturedPiece)`: Reverses a move

### 4. Move (`models/move.go`)

Represents a chess move.

```go
type Move struct {
    From      Position
    To        Position
    Promotion PieceType  // For pawn promotion
}
```

### 5. Game (`game.go`)

Main game engine that orchestrates the game.

**Key Responsibilities**:
- Turn management
- Move validation
- Check/checkmate detection
- Game state management
- Thread-safe operations

**Key Methods**:
- `MakeMove(move)`: Attempts to make a move
- `GetValidMoves(pos)`: Returns all valid moves for a piece
- `IsInCheck()`: Checks if current player is in check
- `GetState()`: Returns current game state

---

## Piece Movement Logic

### Pawn (`pieces/pawn.go`)

**Movement Rules**:
- Moves forward one square (or two from starting position)
- Captures diagonally
- Promotes to Queen when reaching opposite end

**Implementation Highlights**:
- Direction depends on color (White moves up, Black moves down)
- Special handling for starting position (two-square move)
- Promotion logic for end-row positions

### Rook (`pieces/rook.go`)

**Movement Rules**:
- Moves horizontally or vertically any number of squares
- Cannot jump over pieces
- Stops at first piece (captures if opponent, blocked if own)

**Implementation**: Iterates in 4 directions (up, down, left, right) until blocked.

### Knight (`pieces/knight.go`)

**Movement Rules**:
- Moves in L-shape: 2 squares in one direction, 1 square perpendicular
- Can jump over pieces
- 8 possible moves from any position

**Implementation**: Uses fixed offset patterns for all 8 possible moves.

### Bishop (`pieces/bishop.go`)

**Movement Rules**:
- Moves diagonally any number of squares
- Cannot jump over pieces
- Stops at first piece

**Implementation**: Iterates in 4 diagonal directions until blocked.

### Queen (`pieces/queen.go`)

**Movement Rules**:
- Combines Rook and Bishop movement
- Moves horizontally, vertically, or diagonally

**Implementation**: Delegates to both Rook and Bishop behaviors (composition).

### King (`pieces/king.go`)

**Movement Rules**:
- Moves one square in any direction
- Special move: Castling (kingside and queenside)

**Implementation**: 
- 8-directional movement (one square)
- Castling logic (simplified - checks for empty squares and rook presence)

---

## Game Rules Implementation

### Move Validation

1. **Basic Validation**:
   - Source position has a piece
   - Piece belongs to current player
   - Destination doesn't have own piece
   - Move is legal for piece type

2. **Check Prevention**:
   - After making move, checks if own king is in check
   - If check detected, move is rejected and board is restored

### Check Detection

**Algorithm**:
1. Find the king of the given color
2. Check if any opponent piece can attack the king
3. Uses piece behaviors to check attack patterns

### Checkmate Detection

**Algorithm**:
1. Check if king is in check
2. Check if player has any valid moves
3. If in check and no valid moves → Checkmate
4. If not in check but no valid moves → Stalemate

### Valid Move Generation

For each piece:
1. Get all possible moves from piece behavior
2. Try each move on a copy of the board
3. Check if move leaves own king in check
4. Return only moves that don't leave king in check

---

## Usage Examples

### Basic Game Setup

```go
// Create a new game
game := NewGame()

// Print initial board
game.PrintBoard()
```

### Making Moves

```go
// Create move from algebraic notation
from, _ := models.FromAlgebraic("e2")
to, _ := models.FromAlgebraic("e4")
move := models.NewMove(from, to)

// Attempt to make move
if err := game.MakeMove(move); err != nil {
    log.Printf("Move failed: %v", err)
} else {
    fmt.Println("Move successful!")
}
```

### Getting Valid Moves

```go
// Get all valid moves for a piece
pos, _ := models.FromAlgebraic("e2")
validMoves := game.GetValidMoves(pos)

for _, move := range validMoves {
    fmt.Printf("%s ", move.ToAlgebraic())
}
```

### Checking Game State

```go
// Check if current player is in check
if game.IsInCheck() {
    fmt.Println("You are in check!")
}

// Check game state
state := game.GetState()
switch state {
case models.Playing:
    fmt.Println("Game in progress")
case models.WhiteWon:
    fmt.Println("White wins!")
case models.BlackWon:
    fmt.Println("Black wins!")
case models.Stalemate:
    fmt.Println("Stalemate!")
}
```

### Running the Demo

```bash
cd go/chess
go run .
```

---

## Interview Explanation Guide

### How to Explain This Design in an Interview

#### 1. **Start with the Problem Statement**

"I was asked to design a chess game engine. The key requirements were:
- Support all 6 piece types with correct movement rules
- Validate moves and prevent illegal moves
- Detect check, checkmate, and stalemate
- Make it extensible and maintainable
- Ensure thread-safety for concurrent access"

#### 2. **Explain the Core Design Pattern**

"I chose the **Strategy Pattern** as the foundation because:
- Each piece type has different movement rules
- Encapsulates movement logic per piece type
- Makes it easy to modify rules for specific pieces
- Follows Open/Closed Principle - open for extension, closed for modification
- Each piece behavior is independently testable"

**Show the interface**:
```go
type PieceBehavior interface {
    GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move
    CanMove(board *models.Board, from, to models.Position, color models.Color) bool
}
```

#### 3. **Walk Through the Architecture**

**Domain Models**:
- "I created core models: `Position` for board squares, `Piece` for chess pieces, `Board` for the 8x8 grid, and `Move` for representing moves."

**Strategy Pattern**:
- "Each piece type (Pawn, Rook, Knight, etc.) implements the `PieceBehavior` interface with its specific movement logic."

**Game Engine**:
- "The `Game` class orchestrates everything - manages turns, validates moves, detects check/checkmate, and maintains game state."

#### 4. **Explain Move Validation**

"Move validation happens in multiple layers:
1. **Basic validation**: Check if source has piece, piece belongs to player, destination is valid
2. **Piece-specific validation**: Use piece behavior to check if move is legal for that piece type
3. **Check prevention**: After making move, verify it doesn't leave own king in check. If it does, reject and restore board."

#### 5. **Explain Check/Checkmate Detection**

**Check Detection**:
- "Find the king of the given color, then check if any opponent piece can attack it using their movement behaviors."

**Checkmate Detection**:
- "Check if king is in check AND player has no valid moves. If both true → checkmate. If not in check but no valid moves → stalemate."

#### 6. **Highlight Extensibility**

"To add a new piece type (e.g., a custom piece):
1. Create a new struct implementing `PieceBehavior`
2. Implement `GetValidMoves()` and `CanMove()` methods
3. Add it to the `pieceBehaviors` map in Game
4. No changes needed to existing code!

This demonstrates the power of the Strategy Pattern."

#### 7. **Discuss Trade-offs**

**Strategy Pattern**:
- ✅ Clean separation of concerns
- ✅ Easy to modify individual piece rules
- ✅ Testable components
- ❌ More files/classes than monolithic approach

**Move Validation Approach**:
- ✅ Prevents illegal moves
- ✅ Ensures game integrity
- ❌ Requires trying moves (performance consideration)

**Thread-Safety**:
- ✅ Safe for concurrent access
- ❌ Mutex overhead (but necessary for correctness)

#### 8. **Concurrency Considerations**

"I used `sync.RWMutex` for thread-safety:
- Read locks for read operations (`GetBoard()`, `GetState()`)
- Write locks for state-changing operations (`MakeMove()`)
- This allows concurrent reads while preventing race conditions"

---

## Design Decisions

### Why Strategy Pattern?

1. **Separation of Concerns**: Each piece type's movement logic is isolated
2. **Extensibility**: Easy to add new pieces or modify existing ones
3. **Testability**: Each behavior can be unit tested independently
4. **Maintainability**: Changes to one piece don't affect others

### Why Separate Models Package?

1. **Domain Clarity**: Core entities are clearly defined
2. **Reusability**: Models can be used across different components
3. **Testability**: Models can be tested independently
4. **Abstraction**: Game logic doesn't need to know implementation details

### Why Copy Board for Move Validation?

1. **Safety**: Prevents side effects during validation
2. **Correctness**: Can try moves without affecting game state
3. **Simplicity**: Easier to reason about than complex undo/redo logic

### Why Thread-Safe?

1. **Real-world Requirement**: Chess games may have concurrent access (e.g., multiple threads analyzing positions)
2. **Safety**: Prevents race conditions and data corruption
3. **Scalability**: Can handle concurrent move attempts safely

### Why Algebraic Notation?

1. **Standard**: Chess players use algebraic notation (e.g., "e4", "Nf3")
2. **Readability**: More intuitive than numeric coordinates
3. **Compatibility**: Easy to integrate with chess UCI protocols

---

## Extensibility

### Adding a New Piece Type

To add a new piece type (e.g., **Chancellor** - combines Rook and Knight):

#### Step 1: Create the Behavior

```go
// pieces/chancellor.go
package pieces

import "github.com/ishanpandey/LowLevelDesign/go/chess/models"

type ChancellorBehavior struct {
    rook   *RookBehavior
    knight *KnightBehavior
}

func NewChancellorBehavior() *ChancellorBehavior {
    return &ChancellorBehavior{
        rook:   NewRookBehavior(),
        knight: NewKnightBehavior(),
    }
}

func (c *ChancellorBehavior) GetValidMoves(board *models.Board, pos models.Position, color models.Color) []models.Move {
    var moves []models.Move
    moves = append(moves, c.rook.GetValidMoves(board, pos, color)...)
    moves = append(moves, c.knight.GetValidMoves(board, pos, color)...)
    return moves
}

func (c *ChancellorBehavior) CanMove(board *models.Board, from, to models.Position, color models.Color) bool {
    return c.rook.CanMove(board, from, to, color) || c.knight.CanMove(board, from, to, color)
}
```

#### Step 2: Add to Game

```go
// In game.go, add to pieceBehaviors map
game.pieceBehaviors[models.Chancellor] = pieces.NewChancellorBehavior()
```

#### Step 3: Add Piece Type to Models

```go
// In models/piece.go
const (
    // ... existing types
    Chancellor PieceType = iota + 6
)
```

That's it! No other code changes needed.

### Modifying Piece Rules

To modify rules for an existing piece (e.g., make pawns move backward):

1. Open the piece's behavior file (e.g., `pieces/pawn.go`)
2. Modify the `GetValidMoves()` method
3. No changes needed elsewhere!

---

## Key Takeaways for Interviews

1. **Design Pattern**: Strategy Pattern is perfect for different behaviors
2. **SOLID Principles**: 
   - Single Responsibility (each piece behavior does one thing)
   - Open/Closed (open for extension, closed for modification)
   - Dependency Inversion (depend on abstractions, not concretions)
3. **Concurrency**: Thread-safety is crucial for real-world systems
4. **Extensibility**: Design should make adding features easy
5. **Clean Code**: Proper separation of concerns, meaningful names, error handling
6. **Domain Modeling**: Core entities (Board, Piece, Position) are clearly defined

---

## Future Enhancements

Potential improvements for production systems:

1. **En Passant**: Full implementation of en passant capture
2. **Castling Validation**: Check for check, moved pieces, and blocked squares
3. **Move History**: Full move history with undo/redo capability
4. **FEN Support**: Import/export positions in FEN notation
5. **PGN Support**: Parse and generate PGN (Portable Game Notation) files
6. **AI Integration**: Minimax algorithm for computer opponent
7. **Time Controls**: Add time management for timed games
8. **Draw Detection**: Threefold repetition, 50-move rule
9. **Performance**: Optimize move generation (bitboards, move ordering)
10. **UI Integration**: REST API or gRPC for frontend integration

---

## Algorithm Complexity

### Move Generation

- **Time Complexity**: O(n) where n is number of pieces (for check detection)
- **Space Complexity**: O(n) for storing valid moves

### Check Detection

- **Time Complexity**: O(n) where n is number of pieces
- **Space Complexity**: O(1)

### Valid Move Generation

- **Time Complexity**: O(n × m) where n is pieces, m is average moves per piece
- **Space Complexity**: O(m) for storing valid moves

---

## Conclusion

This chess implementation demonstrates:
- ✅ **Clean Architecture**: Proper separation of models, behaviors, and game logic
- ✅ **Design Patterns**: Strategy Pattern for piece behaviors
- ✅ **Concurrency**: Thread-safe operations
- ✅ **Extensibility**: Easy to add new pieces or modify rules
- ✅ **Complete Rules**: Check, checkmate, stalemate detection
- ✅ **Production-Ready**: Error handling, validation, documentation

Perfect for understanding low-level design principles and explaining design decisions in technical interviews!

---

## Author Notes

This implementation focuses on:
1. **Clean Design**: Strategy Pattern for piece behaviors
2. **Complete Rules**: All basic chess rules implemented
3. **Extensibility**: Easy to add features or modify rules
4. **Readability**: Clear code structure for interview explanations

The design emphasizes **readability** and **memorability** to help you explain your thought process clearly in interviews.

