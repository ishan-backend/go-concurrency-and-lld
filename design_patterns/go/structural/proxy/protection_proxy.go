package main

import "fmt"

type Driven interface {
	Drive()
}

type Car struct{}

func (c *Car) Drive() {
	fmt.Println("car is being driven")
}

type Driver struct {
	Age int
}

type VerifyDriverAgeToDriveCarProxy struct {
	car    Car
	driver *Driver
}

func (c *VerifyDriverAgeToDriveCarProxy) Drive() {
	if c.driver.Age > 18 {
		fmt.Println("car is being driven")
	} else {
		fmt.Println("driver is too young")
	}

}

func NewCarProxy(driver *Driver) *VerifyDriverAgeToDriveCarProxy { // introduce dependencies as well before creating proxies
	return &VerifyDriverAgeToDriveCarProxy{Car{}, driver}
}

func main() {
	car := NewCarProxy(&Driver{22})
	car.Drive()
}
