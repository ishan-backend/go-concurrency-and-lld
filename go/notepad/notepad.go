package main

import (
	"fmt"
)

// Notepad represents the main notepad structure
type Notepad struct {
	lines    []string
	clipboard []string
	undoStack []Command
	redoStack []Command
}

// Command interface for undo/redo operations
type Command interface {
	Execute() error
	Undo() error
}

// InsertCommand represents an insert operation
type InsertCommand struct {
	notepad *Notepad
	lineNum int
	text    string
}

func (c *InsertCommand) Execute() error {
	// Adjust line number to 0-based index
	idx := c.lineNum - 1
	if idx < 0 {
		idx = 0
	}
	if idx > len(c.notepad.lines) {
		// Pad with empty lines if needed
		for len(c.notepad.lines) < idx {
			c.notepad.lines = append(c.notepad.lines, "")
		}
	}
	c.notepad.lines = append(c.notepad.lines[:idx], append([]string{c.text}, c.notepad.lines[idx:]...)...)
	return nil
}

func (c *InsertCommand) Undo() error {
	idx := c.lineNum - 1
	if idx < 0 || idx >= len(c.notepad.lines) {
		return fmt.Errorf("invalid line number for undo")
	}
	c.notepad.lines = append(c.notepad.lines[:idx], c.notepad.lines[idx+1:]...)
	return nil
}

// DeleteCommand represents a delete operation
type DeleteCommand struct {
	notepad *Notepad
	startLine int
	endLine   int
	deletedLines []string
}

func (c *DeleteCommand) Execute() error {
	startIdx := c.startLine - 1
	endIdx := c.endLine - 1
	
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx >= len(c.notepad.lines) {
		endIdx = len(c.notepad.lines) - 1
	}
	if startIdx > endIdx || startIdx >= len(c.notepad.lines) {
		return fmt.Errorf("invalid line range")
	}
	
	// Store deleted lines for undo
	c.deletedLines = make([]string, endIdx-startIdx+1)
	copy(c.deletedLines, c.notepad.lines[startIdx:endIdx+1])
	
	// Delete the lines
	c.notepad.lines = append(c.notepad.lines[:startIdx], c.notepad.lines[endIdx+1:]...)
	return nil
}

func (c *DeleteCommand) Undo() error {
	if len(c.deletedLines) == 0 {
		return fmt.Errorf("nothing to undo")
	}
	
	startIdx := c.startLine - 1
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx > len(c.notepad.lines) {
		// Pad with empty lines if needed
		for len(c.notepad.lines) < startIdx {
			c.notepad.lines = append(c.notepad.lines, "")
		}
	}
	
	// Restore deleted lines
	c.notepad.lines = append(c.notepad.lines[:startIdx], append(c.deletedLines, c.notepad.lines[startIdx:]...)...)
	return nil
}

// PasteCommand represents a paste operation
type PasteCommand struct {
	notepad *Notepad
	lineNum int
	pastedLines []string
}

func (c *PasteCommand) Execute() error {
	if len(c.notepad.clipboard) == 0 {
		return fmt.Errorf("clipboard is empty")
	}
	
	idx := c.lineNum - 1
	if idx < 0 {
		idx = 0
	}
	if idx > len(c.notepad.lines) {
		// Pad with empty lines if needed
		for len(c.notepad.lines) < idx {
			c.notepad.lines = append(c.notepad.lines, "")
		}
	}
	
	// Store what was pasted for undo
	c.pastedLines = make([]string, len(c.notepad.clipboard))
	copy(c.pastedLines, c.notepad.clipboard)
	
	// Insert clipboard content
	c.notepad.lines = append(c.notepad.lines[:idx], append(c.notepad.clipboard, c.notepad.lines[idx:]...)...)
	return nil
}

func (c *PasteCommand) Undo() error {
	if len(c.pastedLines) == 0 {
		return fmt.Errorf("nothing to undo")
	}
	
	idx := c.lineNum - 1
	if idx < 0 || idx >= len(c.notepad.lines) {
		return fmt.Errorf("invalid line number for undo")
	}
	
	// Remove pasted lines
	endIdx := idx + len(c.pastedLines)
	if endIdx > len(c.notepad.lines) {
		endIdx = len(c.notepad.lines)
	}
	c.notepad.lines = append(c.notepad.lines[:idx], c.notepad.lines[endIdx:]...)
	return nil
}

// NewNotepad creates a new notepad instance
func NewNotepad() *Notepad {
	return &Notepad{
		lines:     make([]string, 0),
		clipboard: make([]string, 0),
		undoStack: make([]Command, 0),
		redoStack: make([]Command, 0),
	}
}

// Display prints all lines in the notepad
func (n *Notepad) Display() {
	if len(n.lines) == 0 {
		fmt.Println("(empty notepad)")
		return
	}
	for i, line := range n.lines {
		fmt.Printf("%d: %s\n", i+1, line)
	}
}

// Insert inserts text at the specified line number (1-indexed)
func (n *Notepad) Insert(lineNum int, text string) error {
	cmd := &InsertCommand{
		notepad: n,
		lineNum: lineNum,
		text:    text,
	}
	
	if err := cmd.Execute(); err != nil {
		return err
	}
	
	// Clear redo stack when new command is executed
	n.redoStack = make([]Command, 0)
	// Add to undo stack
	n.undoStack = append(n.undoStack, cmd)
	return nil
}

// Delete deletes lines from startLine to endLine (inclusive, 1-indexed)
func (n *Notepad) Delete(startLine, endLine int) error {
	cmd := &DeleteCommand{
		notepad:   n,
		startLine: startLine,
		endLine:   endLine,
	}
	
	if err := cmd.Execute(); err != nil {
		return err
	}
	
	// Clear redo stack when new command is executed
	n.redoStack = make([]Command, 0)
	// Add to undo stack
	n.undoStack = append(n.undoStack, cmd)
	return nil
}

// DeleteSingle deletes a single line (1-indexed)
func (n *Notepad) DeleteSingle(lineNum int) error {
	return n.Delete(lineNum, lineNum)
}

// Undo undoes the last command
func (n *Notepad) Undo() error {
	if len(n.undoStack) == 0 {
		return fmt.Errorf("nothing to undo")
	}
	
	// Pop from undo stack
	cmd := n.undoStack[len(n.undoStack)-1]
	n.undoStack = n.undoStack[:len(n.undoStack)-1]
	
	// Undo the command
	if err := cmd.Undo(); err != nil {
		// Put it back if undo fails
		n.undoStack = append(n.undoStack, cmd)
		return err
	}
	
	// Push to redo stack
	n.redoStack = append(n.redoStack, cmd)
	return nil
}

// Redo redoes the last undone command
func (n *Notepad) Redo() error {
	if len(n.redoStack) == 0 {
		return fmt.Errorf("nothing to redo")
	}
	
	// Pop from redo stack
	cmd := n.redoStack[len(n.redoStack)-1]
	n.redoStack = n.redoStack[:len(n.redoStack)-1]
	
	// Execute the command again
	if err := cmd.Execute(); err != nil {
		// Put it back if execute fails
		n.redoStack = append(n.redoStack, cmd)
		return err
	}
	
	// Push to undo stack
	n.undoStack = append(n.undoStack, cmd)
	return nil
}

// Copy copies lines from startLine to endLine (inclusive, 1-indexed) to clipboard
func (n *Notepad) Copy(startLine, endLine int) error {
	startIdx := startLine - 1
	endIdx := endLine - 1
	
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx >= len(n.lines) {
		endIdx = len(n.lines) - 1
	}
	if startIdx > endIdx || startIdx >= len(n.lines) {
		return fmt.Errorf("invalid line range")
	}
	
	// Copy to clipboard
	n.clipboard = make([]string, endIdx-startIdx+1)
	copy(n.clipboard, n.lines[startIdx:endIdx+1])
	return nil
}

// Paste pastes clipboard content at the specified line number (1-indexed)
func (n *Notepad) Paste(lineNum int) error {
	if len(n.clipboard) == 0 {
		return fmt.Errorf("clipboard is empty")
	}
	
	cmd := &PasteCommand{
		notepad:  n,
		lineNum:  lineNum,
	}
	
	if err := cmd.Execute(); err != nil {
		return err
	}
	
	// Clear redo stack when new command is executed
	n.redoStack = make([]Command, 0)
	// Add to undo stack
	n.undoStack = append(n.undoStack, cmd)
	return nil
}

// Example usage
func main() {
	np := NewNotepad()
	
	fmt.Println("=== Initial State ===")
	np.Display()
	
	fmt.Println("\n=== Inserting lines ===")
	np.Insert(1, "Hello World")
	np.Insert(2, "This is line 2")
	np.Insert(3, "This is line 3")
	np.Display()
	
	fmt.Println("\n=== Copying lines 1-2 ===")
	np.Copy(1, 2)
	fmt.Printf("Clipboard: %v\n", np.clipboard)
	
	fmt.Println("\n=== Pasting at line 4 ===")
	np.Paste(4)
	np.Display()
	
	fmt.Println("\n=== Deleting line 2 ===")
	np.DeleteSingle(2)
	np.Display()
	
	fmt.Println("\n=== Undo (restore line 2) ===")
	np.Undo()
	np.Display()
	
	fmt.Println("\n=== Redo (delete line 2 again) ===")
	np.Redo()
	np.Display()
	
	fmt.Println("\n=== Deleting lines 1-2 ===")
	np.Delete(1, 2)
	np.Display()
	
	fmt.Println("\n=== Undo (restore deleted lines) ===")
	np.Undo()
	np.Display()
}

