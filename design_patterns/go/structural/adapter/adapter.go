package main

import "fmt"

type Line struct {
	x1, x2, y1, y2 int
}

type VectorImage struct {
	lines []Line
}

func NewRectangle(width, height int) *VectorImage {
	width -= 1
	height -= 1
	return &VectorImage{lines: []Line{
		{0, 0, width, 0},
		{0, 0, 0, height},
		{width, 0, width, height},
		{0, height, width, height},
	}}
}

// ^^^ Interface you are given -> You are given lines i.e. []int, but you cannot make them visual on CLI

// Below is the interface you have -> you can work with characters to make rectangle visual

type Point struct {
	X, Y int
}

// RasterImage is our interface which is []Point
type RasterImage interface {
	GetPoints() []Point
}

func DrawPoint(img RasterImage) string {
	// DrawPoint is internal function to print the final rectangle on GUI
	fmt.Println(img)
	return "**********\n**********\n**********\n**********"
}

/*
	adapter is going to be private (not exposed)
	consists of []Point, that is generated from []Line
*/
type vectorToRasterAdapter struct {
	points []Point
}

func (v vectorToRasterAdapter) GetPoints() []Point {
	return v.points
}

func VectorToRaster(vi *VectorImage) RasterImage {
	adapter := vectorToRasterAdapter{}
	for _, line := range vi.lines {
		adapter.AddLine(line)
	}

	return adapter // as RasterImage
}

func (v vectorToRasterAdapter) AddLine(line Line) []Point {
	return []Point{{0, 1}, {1, 1}, {1, 0}, {1, 1}}
}

func main() {
	// Hitting third party NewRectangle API for getting rectangle
	ra := NewRectangle(10, 4)
	adapter := VectorToRaster(ra)
	// Using adapter (new type) and exposing the interface (AddLine); made using simple factory;
	// to use our API DrawPoint
	fmt.Println(DrawPoint(adapter))
}
