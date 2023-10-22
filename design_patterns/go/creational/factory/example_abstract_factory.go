package main

import "fmt"

/*
	iShoe.go: Abstract product
*/
type IShoe interface {
	setLogo(brand string)
	setSize(number int)
	getBrand() string
	getSize() int
}

type Shoe struct {
	logo string
	size int
}

func (s *Shoe) setLogo(brand string) {
	s.logo = brand
}

func (s *Shoe) setSize(number int) {
	s.size = number
}

func (s *Shoe) getBrand() string {
	return s.logo
}

func (s *Shoe) getSize() int {
	return s.size
}

/*
	iShirt.go: Abstract product
*/
type IShirt interface {
	setLogo(brand string)
	setSize(number int)
	getBrand() string
	getSize() int
}

type Shirt struct {
	logo string
	size int
}

func (s *Shirt) setLogo(brand string) {
	s.logo = brand
}

func (s *Shirt) setSize(number int) {
	s.size = number
}

func (s *Shirt) getBrand() string {
	return s.logo
}

func (s *Shirt) getSize() int {
	return s.size
}

/*
	adidasShirt.go: Concrete product
*/
type adidasShirt struct {
	Shirt
}

/*
	nikeShirt.go: Concrete Product
*/
type nikeShirt struct {
	Shirt
}

/*
	adidasShoe.go: Concrete product
*/
type adidasShoe struct {
	Shoe
}

/*
	nikeShoe.go: Concrete Product
*/
type nikeShoe struct {
	Shoe
}

/*
	iSportsFactory.go: Abstract factory interface
*/
type ISportsFactory interface {
	makeShirt() IShirt
	makeShoes() IShoe
	//makeShorts()
}

func getSportsFactory(brandName string) (ISportsFactory, error) {
	if brandName == "adidas" {
		return &Adidas{}, nil
	}
	if brandName == "nike" {
		return &Nike{}, nil
	}

	return nil, fmt.Errorf("wrong brand type passed")
}

/*
	adidas.go: Concrete factory
*/
type Adidas struct{}

func (a *Adidas) makeShoes() IShoe {
	return &adidasShoe{
		Shoe: Shoe{
			logo: "adidas",
			size: 14,
		},
	}
}

func (a *Adidas) makeShirt() IShirt {
	return &adidasShirt{
		Shirt: Shirt{
			logo: "adidas",
			size: 32,
		},
	}
}

/*
	nike.go: Concrete factory
*/
type Nike struct{}

func (n *Nike) makeShoes() IShoe {
	return &nikeShoe{
		Shoe: Shoe{
			logo: "nike",
			size: 14,
		},
	}
}

func (n *Nike) makeShirt() IShirt {
	return &nikeShirt{
		Shirt: Shirt{
			logo: "nike",
			size: 30,
		},
	}
}

/*
	main.go
*/
func main() {
	adidasFactory, _ := getSportsFactory("adidas")
	nikeFactory, _ := getSportsFactory("nike")

	adidasShirt := adidasFactory.makeShirt()
	adidasShoe := adidasFactory.makeShoes()

	nikeShirt := nikeFactory.makeShirt()
	nikeShoe := nikeFactory.makeShoes()

	printShirtDetails(adidasShirt)
	printShoeDetails(adidasShoe)

	printShirtDetails(nikeShirt)
	printShoeDetails(nikeShoe)
}

func printShoeDetails(s IShoe) {
	fmt.Printf(s.getBrand() + fmt.Sprintf("Size: %d", s.getSize()))
	fmt.Println()
}

func printShirtDetails(s IShirt) {
	fmt.Printf(s.getBrand() + fmt.Sprintf("Size: %d", s.getSize()))
	fmt.Println()
}
