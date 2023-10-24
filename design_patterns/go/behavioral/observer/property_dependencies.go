package main

import "fmt"

func (p *Person) IsOkayForInjection() bool {
	if p.GetAge() >= 19 {
		return true
	}
	return false
}

type InjectionManagement struct {
	o Observable
}

func (i *InjectionManagement) Notify(data interface{}) {
	if pc, ok := data.(PropertyChange); ok {
		if pc.Value.(bool) == true {
			fmt.Println("injection given to this patient")

			i.o.UnSubscribe(i) // paracetamol given now doctors(observers) can unsubscribe from this OPD patient
			fmt.Println("okay")
		}
	}
}

func (p *Person) SetAge2(age int) {
	if p.Age == age {
		return
	}

	// cache previous value & compare before sending notification
	oldCanGetInjection := p.IsOkayForInjection()
	p.Age = age
	if oldCanGetInjection != p.IsOkayForInjection() {
		p.SendEventToObserver(PropertyChange{"can get injection", p.IsOkayForInjection()})
	}
}

func main() {
	p := AddPerson("John", 10)
	im := &InjectionManagement{p.Observable}
	p.Subscribe(im) // im is now subscribed to events happening on p

	for i := 10; i < 20; i++ {
		fmt.Println("setting age to ", i)
		p.SetAge2(i)
	}
}
