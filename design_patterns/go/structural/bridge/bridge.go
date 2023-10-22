package main

import "fmt"

type Renderer interface {
	RenderCircle(radius float32)
	//RenderRectangle()
}

type VectorRenderer struct {
}

func (v *VectorRenderer) RenderCircle(radius float32) {
	fmt.Println("we are drawing a circle of radius", radius)
}

type RastorRenderer struct {
	Dpi int
}

func (r *RastorRenderer) RenderCircle(radius float32) {
	fmt.Println("we are drawing pixels for rastor circle of radius", radius)
}

type Circle struct {
	renderer Renderer
	radius   float32
}

func NewCircle(renderer Renderer, radius float32) *Circle {
	return &Circle{
		renderer: renderer,
		radius:   radius,
	}
}

func (c *Circle) Draw() {
	c.renderer.RenderCircle(c.radius)
}

func (c *Circle) Resize(factor float32) {
	c.radius *= factor
}

func main() {
	raster := &RastorRenderer{}
	vector := &VectorRenderer{}

	circle := &Circle{
		renderer: raster,
		radius:   2,
	}
	circle.Draw()
	circle.Resize(0.5)
	circle.Draw()

	circle2 := &Circle{
		renderer: vector,
		radius:   10,
	}
	circle2.Draw()
	circle2.Resize(0.5)
	circle2.Draw()
}

/*
	When you introduce new shape - say square

    type Square struct {
		renderer Renderer
		side   float32
	}

	func NewSquare(renderer Renderer, side float32) *Square {
		return &Circle{
			renderer: renderer,
			radius:   radius,
		}
	}

	func (c *Square) Draw() {
		c.renderer.RenderCircle(c.side)
	}

	func (c *Square) Resize(factor float32) {
		c.radius *= factor
	}

	Also you would need to add a method in interface Renderer and subsequent downstream implementations from VectorRenderer & RastorRenderer

*/
