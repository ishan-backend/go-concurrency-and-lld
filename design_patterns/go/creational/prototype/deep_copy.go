package main

import "fmt"

type Address struct {
	AddressLine1 string
	AddressLine2 string
	Tehsil       string
	Gram         string
	District     string
	State        string
	PinCode      string
}

func (a *Address) DeepCopy() *Address {
	return &Address{
		AddressLine1: a.AddressLine1,
		District:     a.District,
		State:        a.State,
	}
}

type Person struct {
	Name    string
	Address *Address
	Friends []string
}

func (p *Person) DeepCopy() *Person {
	q := *p // Using this you are making copy of actual object passed, and not the copy of pointer --> dont use q := p
	q.Address = p.Address.DeepCopy()
	copy(q.Friends, p.Friends) // destination, source
	return &q
}

func main() {
	ram := Person{Name: "Ram Ji", Address: &Address{AddressLine1: "Janmsthan", District: "Ayodhya Ji", State: "Uttar Pradesh"}}

	/*
		// This copies ram object, ram.Address is a pointer and gets copied too, so now if you change shyam.Address, ram.Address will change too
		// This process involves copying the pointer
		// Modifying either of original and copy, would result in modification

		shyam := ram
		shyam.Name = "Shyam" // ok
		shyam.Address = &Address{AddressLine1: "Janmsthan", District: "Mathura", State: "Uttar Pradesh"}

		fmt.Println(shyam) // will print name, and address of pointer
		fmt.Println(shyam, shyam.Address)
		fmt.Println(ram, ram.Address)
	*/

	/*
			// deep copy - while copying object, copy everything that it refers to - slices (also a pointer), pointers etc

		shyam := ram
		shyam.Name = "shyam"
		shyam.Address = &Address{
			AddressLine1: ram.Address.AddressLine1,
			District:     ram.Address.District,
			State:        ram.Address.State,
		}
		shyam.Address.District = "Mathura"
		fmt.Println(shyam, shyam.Address)
		fmt.Println(ram, ram.Address)
			// This approach is not scalable, for e.g. with objects having deep recursive structure
	*/

	/*
		// Using Deep Copy function for every type i.e. Copy Method:
		// This is not fine, since we have to double/triple check our structs to have all of its members like structs/slice have a deep copy method available
	*/

	shyam := ram.DeepCopy()
	shyam.Name = "Shyam"
	shyam.Address.District = "Mathura"
	shyam.Friends = append(shyam.Friends, "Radha Rani")
	fmt.Println(shyam, shyam.Address, shyam.Friends)
	fmt.Println(ram, ram.Address, ram.Friends)
}
