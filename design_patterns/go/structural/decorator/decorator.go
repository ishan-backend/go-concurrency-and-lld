package main

import "fmt"

type Shape interface {
	Render() string
}

type Circle struct {
	Radius int
}

func (c *Circle) Render() string {
	return fmt.Sprintf("circle rendered with radius %d", c.Radius)
}

func (c *Circle) Resize(rad int) {
	c.Radius = rad
}

type Square struct {
	Side int
}

func (s *Square) Render() string {
	return fmt.Sprintf("square rendered with side %d", s.Side)
}

/*
	How add a functionality to add colors to your Square & Circle?
	- If you want to add a attribute to struct, that will break OCP.
		Because, some logic might be sitting already - say serialisation of struct that would break
	- One option: Aggregation using new struct using generics , which is not available in Go
	-
*/

type ColoredShape struct { // decorator
	Shape Shape // shape its decorating
	Color string
}

func (c *ColoredShape) Render() string {
	return fmt.Sprintf("%s has color %s", c.Shape.Render(), c.Color)
}

/*
	Composition of decorators
*/

type TransparentShape struct {
	Shape        Shape
	Transparency float32
}

func (c *TransparentShape) Render() string {
	return fmt.Sprintf("%s with transparency %f", c.Shape.Render(), c.Transparency)
}

func main() {
	circle := Circle{2}
	fmt.Println(circle.Render())

	redCircle := ColoredShape{ // calling decorator
		Shape: &Circle{5},
		Color: "red",
	}
	fmt.Println(redCircle.Render())

	// Composition of decorators - either way
	rts := TransparentShape{&redCircle, 0.1}
	fmt.Println(rts.Render())

	trs := ColoredShape{&TransparentShape{
		Shape:        &Circle{3},
		Transparency: 0.1,
	}, "red"}
	fmt.Println(trs.Render())
}
