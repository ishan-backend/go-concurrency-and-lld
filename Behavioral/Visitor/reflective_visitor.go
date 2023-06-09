package main

import (
	"fmt"
	"strings"
)

// (1+2)+3
type Expression interface{}

type DoubleExpression struct {
	value float64
}

type AdditionExpression struct {
	left, right Expression
}

// Print takes an expression and prints it to sb
// Separation of concerns
func Print(e Expression, sb *strings.Builder) {
	if de, ok := e.(*DoubleExpression); ok {
		sb.WriteString(fmt.Sprintf("%g", de.value))
	} else if ae, ok := e.(*AdditionExpression); ok {
		sb.WriteRune('(')
		Print(ae.left, sb)
		sb.WriteRune('+')
		Print(ae.right, sb)
		sb.WriteRune(')')
	}

	// Problematic: If a new type is added, you have to add a new code here to handle that type
	// We are breaking open-closed principle / OCP. We don't want to touch existing code that is tested.
}

func main() {
	e := &AdditionExpression{
		left: &AdditionExpression{
			left:  &DoubleExpression{value: 1},
			right: &DoubleExpression{value: 2},
		},
		right: &DoubleExpression{value: 3},
	}

	sb := strings.Builder{}
	Print(e, &sb)
	fmt.Println(sb.String())
}
