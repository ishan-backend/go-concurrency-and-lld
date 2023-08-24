package main

import "fmt"

// Create Factories for specific roles within a company
type Employee struct {
	Name, Position string
	AnnualIncome   int
}

// functional approach - [High order function]factory function returns a function, that you can subsequently use to fully initialize an object personally.
func NewEmployeeFactory(position string, annualIncome int) func(name string) *Employee {
	return func(name string) *Employee {
		return &Employee{name, position, annualIncome}
	}
}

// structural approach - makes factory a struct
type EmployeeFactory struct {
	Position     string
	AnnualIncome int
}

func NewEmployeeFactory2(position string, annualIncome int) *EmployeeFactory {
	return &EmployeeFactory{position, annualIncome}
}

func (e *EmployeeFactory) Create(name string) *Employee {
	return &Employee{name, e.Position, e.AnnualIncome}
}

func main() {
	// Functional approach
	/*
		// factory functions
		sde1Factory := NewEmployeeFactory("Developer", 2400000)
		emFactory := NewEmployeeFactory("EM", 10000000)

		// Invoke factory function to create a developer
		sde1 := sde1Factory("Punit")
		em := emFactory("Akhil Gupta")

		fmt.Println(sde1, em)
	*/

	// Structural approach
	sde1Factory := NewEmployeeFactory2("Developer", 2400000)
	emFactory := NewEmployeeFactory2("EM", 100000000)

	sde1Factory.AnnualIncome = 3000000
	sde1 := sde1Factory.Create("Punit")
	em := emFactory.Create("Akhil Gupta")

	fmt.Println(sde1, em)
}

/*

In passing ordinary functions into something is easier than passing in a specialized object, because

for a specialized object, whoever is consuming that object has to explicitly know that there is some

sort of create method and that they have to call this create method.



So this is a situation where, for example, you might try to introduce some sort of interface which

tells you explicitly that there is a create method there.

And here are the arguments.

And then you could also use this kind of interface approach to to pass an interface of the factory rather

than the factory itself.

So that's a possibility.

But whichever option you go, both of them are just fine.

The first one is probably more idiomatic, more kind of more functional, more idiomatic.

So I would recommend the first option.

*/
