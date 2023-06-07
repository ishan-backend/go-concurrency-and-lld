package main

import "fmt"

type Creature struct {
	Name            string
	Attack, Defense int
}

func NewCreature(name string, attack, defense int) *Creature {
	return &Creature{Name: name, Attack: attack, Defense: defense}
}

func (c *Creature) String() string {
	return fmt.Sprintf("%s, %v, %v", c.Name, c.Attack, c.Defense)
}

/*
	When you add a modifier, you add it on top of an existing modifier
*/
type Modifier interface {
	Add(m Modifier)
	Apply()
}

// is effectively singly linked list, and implements Modifier interface
type CreatureModifier struct {
	creature *Creature
	next     Modifier
}

func (c *CreatureModifier) Add(m Modifier) {
	if c.next != nil {
		c.next.Add(m)
	} else {
		c.next = m
	}
}

func (c *CreatureModifier) Apply() {

}

func main() {

}
