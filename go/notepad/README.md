# Notepad Implementation

A simple notepad implementation in Go with clipboard, undo, and redo functionality.

## Features

- **display()** - Display all lines in the notepad
- **insert(n, text)** - Insert text at line n (1-indexed)
- **delete(n, m)** - Delete lines from n to m (inclusive, 1-indexed)
- **delete(n)** - Delete single line n (1-indexed)
- **undo()** - Undo the last command
- **redo()** - Redo the last undone command
- **copy(n, m)** - Copy lines n to m to clipboard
- **paste(n)** - Paste clipboard content at line n

## Key Implementation Details

### Clipboard
- Stores copied lines in a slice
- Persists until next copy operation
- Can be pasted multiple times

### Undo/Redo
- Uses Command Pattern with two stacks:
  - **Undo Stack**: Stores executed commands
  - **Redo Stack**: Stores undone commands
- Each command implements `Execute()` and `Undo()` methods
- New commands clear the redo stack
- Supports undo/redo for: Insert, Delete, and Paste operations

## Usage

```go
np := NewNotepad()
np.Insert(1, "Hello World")
np.Insert(2, "Line 2")
np.Copy(1, 2)
np.Paste(3)
np.Display()
np.Undo()
np.Redo()
```

## Run

```bash
go run notepad.go
```

Analyzing the notepad code to identify all API flows and summarize what each does.


[2 tools called]


## Notepad API Flows

### 1. **NewNotepad()**
Creates a new notepad with empty lines, clipboard, undo stack, and redo stack.

### 2. **Display()**
Prints all lines with 1-indexed line numbers. Shows "(empty notepad)" if there are no lines.

### 3. **Insert(lineNum, text)**
Inserts text at the given line number (1-indexed). If the line number is beyond the end, pads with empty lines. Creates an `InsertCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 4. **Delete(startLine, endLine)**
Deletes lines from `startLine` to `endLine` (inclusive, 1-indexed). Validates the range, stores deleted lines for undo, removes them, creates a `DeleteCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 5. **DeleteSingle(lineNum)**
Deletes a single line by calling `Delete(lineNum, lineNum)`.

### 6. **Copy(startLine, endLine)**
Copies lines from `startLine` to `endLine` (inclusive, 1-indexed) to the clipboard. Validates the range, replaces the clipboard with the selected lines, and does not affect undo/redo.

### 7. **Paste(lineNum)**
Pastes clipboard content at the given line number (1-indexed). If the clipboard is empty, returns an error. If the line number is beyond the end, pads with empty lines. Creates a `PasteCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 8. **Undo()**
Reverses the last executed command. Pops from the undo stack, calls `Undo()` on the command, and pushes it to the redo stack. If undo fails, restores the command to the undo stack. Returns an error if there's nothing to undo.

### 9. **Redo()**
Re-executes the last undone command. Pops from the redo stack, calls `Execute()` on the command, and pushes it to the undo stack. If execution fails, restores the command to the redo stack. Returns an error if there's nothing to redo.

## Design patterns

- Command pattern: Insert, Delete, and Paste are commands with `Execute()` and `Undo()`.
- Undo/redo stacks: Executed commands go to the undo stack; undone commands go to the redo stack.
- Clipboard: Stores copied lines until the next copy operation.

## Notes

- Line numbers are 1-indexed in the API but converted to 0-indexed internally.
- New commands clear the redo stack.
- Copy does not create a command and does not affect undo/redo.
- Insert, Delete, and Paste are undoable/redoable.




-------------------------------------------------------


Analyzing the notepad code to identify all API flows and summarize what each does.


[2 tools called]


## Notepad API Flows

### 1. **NewNotepad()**
Creates a new notepad with empty lines, clipboard, undo stack, and redo stack.

### 2. **Display()**
Prints all lines with 1-indexed line numbers. Shows "(empty notepad)" if there are no lines.

### 3. **Insert(lineNum, text)**
Inserts text at the given line number (1-indexed). If the line number is beyond the end, pads with empty lines. Creates an `InsertCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 4. **Delete(startLine, endLine)**
Deletes lines from `startLine` to `endLine` (inclusive, 1-indexed). Validates the range, stores deleted lines for undo, removes them, creates a `DeleteCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 5. **DeleteSingle(lineNum)**
Deletes a single line by calling `Delete(lineNum, lineNum)`.

### 6. **Copy(startLine, endLine)**
Copies lines from `startLine` to `endLine` (inclusive, 1-indexed) to the clipboard. Validates the range, replaces the clipboard with the selected lines, and does not affect undo/redo.

### 7. **Paste(lineNum)**
Pastes clipboard content at the given line number (1-indexed). If the clipboard is empty, returns an error. If the line number is beyond the end, pads with empty lines. Creates a `PasteCommand`, executes it, clears the redo stack, and pushes the command to the undo stack.

### 8. **Undo()**
Reverses the last executed command. Pops from the undo stack, calls `Undo()` on the command, and pushes it to the redo stack. If undo fails, restores the command to the undo stack. Returns an error if there's nothing to undo.

### 9. **Redo()**
Re-executes the last undone command. Pops from the redo stack, calls `Execute()` on the command, and pushes it to the undo stack. If execution fails, restores the command to the redo stack. Returns an error if there's nothing to redo.

## Design patterns

- Command pattern: Insert, Delete, and Paste are commands with `Execute()` and `Undo()`.
- Undo/redo stacks: Executed commands go to the undo stack; undone commands go to the redo stack.
- Clipboard: Stores copied lines until the next copy operation.

## Notes

- Line numbers are 1-indexed in the API but converted to 0-indexed internally.
- New commands clear the redo stack.
- Copy does not create a command and does not affect undo/redo.
- Insert, Delete, and Paste are undoable/redoable.