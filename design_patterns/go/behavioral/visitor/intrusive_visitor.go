package main

import (
	"fmt"
	"strings"
)

// (1+2)+3
type Expression interface {
	// Violation of open-closed principle
	Print(sb *strings.Builder)
}

type DoubleExpression struct {
	value float64
}

func (db *DoubleExpression) Print(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%g", db.value))
}

type AdditionExpression struct {
	left, right Expression
}

// (left + right)
func (ab *AdditionExpression) Print(sb *strings.Builder) {
	sb.WriteRune('(')
	ab.left.Print(sb)
	sb.WriteRune('+')
	ab.right.Print(sb)
	sb.WriteRune(')')
}

func main() {
	e := &AdditionExpression{
		left: &AdditionExpression{
			left:  &DoubleExpression{value: 1},
			right: &DoubleExpression{value: 2},
		},
		right: &DoubleExpression{value: 3},
	}

	// Intrusion visitor violates open-closed principle, as Expression interface originally defined empty interface, needs a method Print() to be added
	// for Printing the expression required from given hierarchy
	sb := strings.Builder{}
	e.Print(&sb)
	fmt.Println(sb.String())

	// string builder gets passed here to various types - visitor
	// this is intrusive, -> you modify interface part of element hierarchy as well as elements themselves
	// in this setup, you would have to add Evaluate here, to calculate the sum of this expression
	// Seperation of concerns & Single Responsibility Principle is not fulfilled here.
}
