package main

import (
	"fmt"
)

// In factory function implementation, you can return interface that this struct/whole-sale object conforms to.
// You can use the interface to modify the underlying type.
type Person interface {
	SayHello()
}

type person struct {
	name string
	age  int
}

func (p *person) SayHello() {
	fmt.Printf("name: %s, age: %v ", p.name, p.age)
}

type tiredPerson struct {
	name string
	age  int
}

func (p *tiredPerson) SayHello() {
	fmt.Printf("name: %s, age: %v but implementation is different", p.name, p.age)
}

// When return type is interface, you don't have to put *
// You can have different types of objects returned from interface
func NewPerson(name string, age int) Person {
	if age > 60 {
		return &tiredPerson{name: name, age: age}
	}
	return &person{name: name, age: age}
}

func main() {
	p := NewPerson("Ishan", 22)
	p.SayHello()
	// Above is neat way of encapsulating information
	// you can't do:  p.age++

}
