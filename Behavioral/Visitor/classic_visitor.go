package main

import (
	"fmt"
	"strings"
)

/*
// (1+2)+3
type Expression interface{}

type DoubleExpression struct {
	value float64
}

type AdditionExpression struct {
	left, right Expression
}

// Function overloading not allowed
func Print(e AdditionExpression, sb *strings.Builder) {
	sb.WriteRune('(')
	Print(e.left, sb)
	sb.WriteRune('+')
	Print(e.right, sb) // Expression (interface) -> cannot identify which print to call at compile time
	sb.WriteRune(')')
}

func Print(e DoubleExpression, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%g", de.value))
}
*/

/*
	Double Dispatch: = double jump, Accept()->Visit
	- being able to call correct method, not just on the basis of argument but also on the basis of who the caller is

	- You can modify the Expression interface only once, with double dispatch
*/

type ExpressionVisitor interface {
	VisitDoubleExpression(e *DoubleExpression)
	VisitAdditionExpression(e *AdditionExpression)

	// You cannot forget to handle any type of new expression you want to add e.g. - ,  *
}

type Expression interface {
	Accept(ev ExpressionVisitor)
}

type DoubleExpression struct {
	value float64
}

func (d *DoubleExpression) Accept(ev ExpressionVisitor) {
	ev.VisitDoubleExpression(d)
}

type AdditionExpression struct {
	left, right Expression
}

func (a *AdditionExpression) Accept(ev ExpressionVisitor) {
	ev.VisitAdditionExpression(a)
}

type ExpressionPrinter struct { // This is visitor, and its job is to visit and Print i.e. a functionality ; like thus for any other functionality create a new visitor e.g. Evaluate
	sb strings.Builder
}

func (e *ExpressionPrinter) VisitDoubleExpression(d *DoubleExpression) {
	e.sb.WriteString(fmt.Sprintf("%g", d.value))
}

func (e *ExpressionPrinter) VisitAdditionExpression(a *AdditionExpression) {
	e.sb.WriteRune('(')
	a.left.Accept(e) // Pass the visitor to Accept
	e.sb.WriteRune('+')
	a.right.Accept(e)
	e.sb.WriteRune(')')
}

func NewExpressionPrinter() *ExpressionPrinter {
	return &ExpressionPrinter{sb: strings.Builder{}}
}

func (ep *ExpressionPrinter) String() string {
	return ep.sb.String()
}

func main() {
	e := &AdditionExpression{
		left: &AdditionExpression{
			left:  &DoubleExpression{value: 1},
			right: &DoubleExpression{value: 2},
		},
		right: &DoubleExpression{value: 3},
	}

	ep := NewExpressionPrinter()
	e.Accept(ep)
	fmt.Println(ep.String())
}
