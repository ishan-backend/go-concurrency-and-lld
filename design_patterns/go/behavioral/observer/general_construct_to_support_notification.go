package main // main

import (
	"container/list"
	"fmt"
)

type Observable struct {
	subs *list.List // list of subscribers - which want to listen to events happening on Observable
}

func (o *Observable) Subscribe(x Observer) {
	// Subscribe adds x to list of subs for this observable
	o.subs.PushBack(x)
}

func (o *Observable) UnSubscribe(x Observer) {
	// UnSubscribe removes an observer from list of subs for this observable
	for z := o.subs.Front(); z != nil; z = z.Next() {
		if z.Value.(Observer) == x {
			o.subs.Remove(z)
		}
	}
}

func (o *Observable) SendEventToObserver(data interface{}) {
	// SendEventToObserver will actually notify all the observers that something happened to observable
	for z := o.subs.Front(); z != nil; z = z.Next() {
		z.Value.(Observer).Notify(data)
	}
}

// Observer is an interface, which gets notified when something happens to observable, any struct that implements Observer will have this type
type Observer interface {
	Notify(data interface{}) // you would need to send some service data about what happened/changed to x
}

// --------------------------------------------------------------------------------------------------------------

// Person once gets sick, and wants to send out notification for some action (for any doctor to be assigned)
type Person struct {
	Observable
	Name string
	Age  int
}

func AddPerson(name string, age int) Person {
	return Person{
		Observable: Observable{new(list.List)},
		Name:       name,
		Age:        age,
	}
}

func (p *Person) BecomeSick() {
	p.SendEventToObserver(fmt.Sprintf(p.Name + " have become sick"))
}

// DoctorService lets say is one such service listening to notifications
type DoctorService struct{}

func (d *DoctorService) Notify(data interface{}) {
	fmt.Printf("A doctor has been called for %s", data.(string))
}

//func main() {
//	p := AddPerson("Jon Snow", 18)
//	ds := &DoctorService{}
//	p.Subscribe(ds)
//
//	p.BecomeSick()
//}
