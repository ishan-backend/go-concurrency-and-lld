package main

import (
	"fmt"
	"strings"
)

// GraphicObject - to console/print it, is a recursive relationship
// GraphicObject can contain infinite number of GraphicObject(s)
type GraphicObject struct {
	Name, Color string

	// one figure consists of a bunch of shapes
	Children []GraphicObject
}

func (g *GraphicObject) String() string {
	sb := strings.Builder{}
	g.print(&sb, 0)
	return sb.String()
}

// print is an algorithm - which do not really care about children is nil or not. Into the depth of object as far as necessary
func (g *GraphicObject) print(sb *strings.Builder, depth int) {
	sb.WriteString(strings.Repeat("*", depth))
	if len(g.Color) > 0 {
		sb.WriteString(g.Color)
		sb.WriteRune(' ')
	}
	sb.WriteString(g.Name)
	sb.WriteRune(' ')
	sb.WriteString("\n")
	for _, child := range g.Children {
		child.print(sb, depth+1)
	}
}

/*
	Concrete Types
*/
func NewCircle(color string) *GraphicObject {
	return &GraphicObject{
		Name:     "Circle",
		Color:    color,
		Children: nil,
	}
}

func NewSquare(color string) *GraphicObject {
	return &GraphicObject{
		Name:     "Square",
		Color:    color,
		Children: nil,
	}
}

func main() {
	drawing := GraphicObject{"My Drawing", "red", nil}
	drawing.Children = append(drawing.Children, *NewCircle("yellow"))
	drawing.Children = append(drawing.Children, *NewSquare("blue"))

	group := GraphicObject{"Group 1", "", nil}
	group.Children = append(group.Children, *NewCircle("yellow"))
	group.Children = append(group.Children, *NewSquare("blue"))

	drawing.Children = append(drawing.Children, group)
	fmt.Println(drawing.String())
}
