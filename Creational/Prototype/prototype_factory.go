package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

/*
	Employees and Addresses

	- Employees can work from different locations
	- Employees can have different designations

	- Pre-defined Address for the main office
	- Pre-defined Address for auxiliary office


	Challenge: Too much customisations from a set of functions

	Solution: Have functions which allow you to create people working in different offices.
*/

type Address struct {
	FlatNo                  int
	StreetName, City, State string
}

type Employee struct {
	Name    string
	Address Address
}

func (p *Employee) DeepCopy() *Employee {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	_ = e.Encode(p)

	fmt.Println(string(b.Bytes()))

	d := gob.NewDecoder(&b)
	result := Employee{}
	_ = d.Decode(&result)
	return &result
}

// ** To implement a prototype, partially construct an object and store it somewhere **
// We want to make a copy of the prototype and then allow the person to customize it
var mainOfficeEmployeePrototype = Employee{"", Address{0, "Trinity Church", "Bangalore", "Karnataka"}}
var auxOfficeEmployeePrototype = Employee{"", Address{0, "Thane Road", "Mumbai", "Maharashtra"}}

// Utility Function: Takes a protoType employee and returns the customized employee
func newEmployee(proto *Employee, name string, flatNo int) *Employee {
	// ** Deep copy the prototype **
	result := proto.DeepCopy()
	// ** Customize the resulting instance**
	result.Name = name
	result.Address.FlatNo = flatNo
	return result
}

// ** Conveinient APIs for using prototypes **
func newMainOfficeEmployee(name string, flatNo int) *Employee {
	return newEmployee(&mainOfficeEmployeePrototype, name, flatNo)
}
func newAuxOfficeEmployee(name string, flatNo int) *Employee {
	return newEmployee(&auxOfficeEmployeePrototype, name, flatNo)
}

// Prototype Factory is merely a conveinence approach
func main() {
	john := newMainOfficeEmployee("John", 1)
	jane := newAuxOfficeEmployee("Jane", 2)

	fmt.Println(john, jane)
}
