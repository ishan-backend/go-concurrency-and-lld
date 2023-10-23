package main

import "fmt"

type Lion struct {
	power int
}

func (l *Lion) Power() int {
	return l.power
}

func (l *Lion) SetPower(power int) {
	l.power = power
}

func (l *Lion) Attack(power int) {
	if l.Power() >= 100 {
		fmt.Println("Lion attacks with power ", power)
	}
}

type Man struct {
	power int
}

func (l *Man) Power() int {
	return l.power
}

func (l *Man) SetPower(power int) {
	l.power = power
}

func (l *Man) Run(power int) {
	if l.Power() >= 20 {
		fmt.Println("Man is running with power ", power)
	}
}

/*
	Construction of Vishnu would not be straight-forward aggregation
	- keep aggregation via fields
	- redefine the behaviours Attack() and Run(); basically redirect/proxy them
*/

type Vishnu struct {
	lion  Lion
	human Man
}

func (v *Vishnu) Power() int {
	return v.lion.Power()
}

func (v *Vishnu) SetPower(power int) {
	v.lion.SetPower(power) // redirect / proxy
	v.human.SetPower(power)
}

func (v *Vishnu) Attack() {
	v.lion.Attack(v.Power())
}

func (v *Vishnu) Run() {
	v.human.Run(v.Power())
}

func main() {
	v := Vishnu{}
	v.SetPower(100)
	v.Run()
	v.Attack()
}
