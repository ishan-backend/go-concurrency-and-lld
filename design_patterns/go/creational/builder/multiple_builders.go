package main

import "fmt"

type Person struct {
	// address info
	StreetName, City, State string
	Pincode                 int

	// job info
	CompanyName, Position string
	AnnualIncome          int
}

// PersonBuilder is starting point for building up a person
type PersonBuilder struct {
	p *Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{p: &Person{}}
}

/*
	Utility methods on PersonBuilder that give us PersonAddressBuilder, PersonJobBuilder

	These methods give you ability to jump to different builders, since internally it's the same PersonBuilder pointer
*/
func (p *PersonBuilder) Lives() *PersonAddressBuilder {
	return &PersonAddressBuilder{*p}
}

type PersonAddressBuilder struct {
	PersonBuilder
}

func (p *PersonBuilder) Works() *PersonJobBuilder {
	return &PersonJobBuilder{*p}
}

type PersonJobBuilder struct {
	PersonBuilder
}

/*
	Methods on PersonAddressBuilder & PersonJobBuilder methods to populate the values to PersonBuilder
*/

func (a *PersonAddressBuilder) At(streetName string) *PersonAddressBuilder {
	a.p.StreetName = streetName
	return a
}

func (a *PersonAddressBuilder) In(city string) *PersonAddressBuilder {
	a.p.City = city
	return a
}

func (a *PersonAddressBuilder) WithPinCode(pinCode int) *PersonAddressBuilder {
	a.p.Pincode = pinCode
	return a
}

func (j *PersonJobBuilder) As(position string) *PersonJobBuilder {
	j.p.Position = position
	return j
}

func (j *PersonJobBuilder) In(companyName string) *PersonJobBuilder {
	j.p.CompanyName = companyName
	return j
}

func (j *PersonJobBuilder) WithIncome(annualIncome int) *PersonJobBuilder {
	j.p.AnnualIncome = annualIncome
	return j
}

// Build is the method to be called to return final Person object once all inits are done
func (b *PersonBuilder) Build() *Person {
	return b.p
}

func main() {
	pb := NewPersonBuilder()
	pb.Lives().At("123 London Road").In("London").WithPinCode(340089).Works().As("Chairman").In("Pratt & Whitney").WithIncome(1000000000)

	person := pb.Build()
	fmt.Println(person)
}
