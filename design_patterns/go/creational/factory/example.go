package main

import "fmt"

/*
	iGun.go - Product interface
*/
type IGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

/*
	gun.go - Concrete Product
*/
type gun struct {
	name  string
	power int
}

func (g *gun) setName(name string) {
	g.name = name
}

func (g *gun) setPower(power int) {
	g.power = power
}

func (g *gun) getName() string {
	return g.name
}

func (g *gun) getPower() int {
	return g.power
}

/*
	ak47.go - Concrete product
*/
type AK47 struct {
	gun
}

func newAK47() IGun {
	return &AK47{
		gun: gun{
			name:  "AK47 Rifle",
			power: 400,
		}}
}

/*
	musket.go - Concrete product
*/
type musket struct {
	gun
}

func newMusket() IGun {
	return &musket{
		gun: gun{
			name:  "Musket",
			power: 200,
		}}
}

/*
	gunFactory.go - Factory
*/
func createGun(gunType string) (IGun, error) {
	if gunType == "ak47" {
		return newAK47(), nil
	} else if gunType == "musket" {
		return newMusket(), nil
	}
	return nil, fmt.Errorf("wrong gun type passed")
}

/*
	main.go
*/
func main() {
	ak47, _ := createGun("ak47")
	musket, _ := createGun("musket")

	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g IGun) {
	fmt.Printf("Gun: %s", g.getName())
	fmt.Println()
	fmt.Printf("Power: %d", g.getPower())
	fmt.Println()
}
