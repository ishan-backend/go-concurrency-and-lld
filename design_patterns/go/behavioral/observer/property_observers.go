package main

import (
	"fmt"
)

/*
	Old Code
*/

/*
	New Code
*/

type PropertyChange struct {
	Name  string // name of property from subject
	Value interface{}
}

func (p *Person) GetAge() int {
	return p.Age
}

func (p *Person) SetAge(age int) { // setter implementation is important
	if p.Age == age {
		return
	}
	p.Age = age
	p.SendEventToObserver(PropertyChange{"age", p.Age})
}

type MedicineManagement struct {
	o Observable
}

func (m *MedicineManagement) Notify(data interface{}) {
	if pc, ok := data.(PropertyChange); ok {
		if pc.Value.(int) >= 18 {
			fmt.Println("paracetamol can be given to patient")
			m.o.UnSubscribe(m) // paracetamol given now doctors(observers) can unsubscribe from this OPD patient
		}
	}
}

func main2() {
	p := AddPerson("John", 20)
	t := &MedicineManagement{p.Observable}
	p.Subscribe(t)

	for i := 16; i <= 20; i++ {
		fmt.Println("setting age to value ", i)
		p.SetAge(i)
	}
}
