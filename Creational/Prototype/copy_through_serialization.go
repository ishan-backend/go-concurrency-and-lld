package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// You can use either of binary serialisation / json serialisation, we are using binary serialisation
type Address struct {
	AddressLine1 string
	AddressLine2 string
	Tehsil       string
	Gram         string
	District     string
	State        string
	PinCode      string
}

type Person struct {
	Name    string
	Address *Address
	Friends []string
}

func (p *Person) DeepCopy() *Person {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	_ = e.Encode(p)

	fmt.Println(string(b.Bytes()))

	d := gob.NewDecoder(&b)
	result := Person{}
	_ = d.Decode(&result)
	return &result
}

func main() {
	ram := Person{Name: "Ram Ji", Address: &Address{AddressLine1: "Janmsthan", District: "Ayodhya Ji", State: "Uttar Pradesh"}}

	shyam := ram.DeepCopy()
	shyam.Name = "Shyam"
	shyam.Address.District = "Mathura"
	shyam.Friends = append(shyam.Friends, "Radha Rani")
	fmt.Println(shyam, shyam.Address, shyam.Friends)
	fmt.Println(ram, ram.Address, ram.Friends)
}
