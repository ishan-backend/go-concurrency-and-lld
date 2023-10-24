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
	Abstract Modifier
	When you add a modifier, you add it on top of an existing modifier
	And we want to build a stack/list of these modifications, so that these can be applied to creature one after other
*/
type Modifier interface {
	Add(m Modifier) // adds additional modifier to chain of responsibility, attaching it to end of already applied modifiers
	Apply()         // modifier gets applied when this method is called
}

// CreatureModifier concrete type - what creature this Modifier is applied to
type CreatureModifier struct {
	creature *Creature
	next     Modifier
}

func (c *CreatureModifier) Add(m Modifier) {
	if c.next != nil {
		c.next.Add(m) // add element to the end
	} else {
		c.next = m // replace a modifier
	}
}

// Apply is just calling other modifiers on this, and is actually not doing anything by itself.
// Reason is: we are going to aggregate it and make use of it
func (c *CreatureModifier) Apply() {
	if c.next != nil {
		c.next.Apply()
	}
}

func NewCreatureModifier(creature *Creature) *CreatureModifier { // Factory for CreatureModifier
	return &CreatureModifier{creature: creature}
}

/*
	Concrete modifier definitions
	We are initialising the aggregated part for creature modifier
	We first initialise
*/
type DoubleAttackModifier struct {
	CreatureModifier
}

func NewDoubleAttackModifier(c *Creature) *DoubleAttackModifier {
	return &DoubleAttackModifier{
		CreatureModifier{creature: c},
	}
}

func (d *DoubleAttackModifier) Apply() {
	fmt.Println("doubling ", d.creature.Name, " attack value")
	d.creature.Attack *= 2
	d.CreatureModifier.Apply() // => this will call (c *CreatureModifier) Apply()
}

/*
	Increase Defense Modifier
*/
type IncreaseDefenseModifier struct {
	CreatureModifier
}

func NewIncreasedDefenseModifier(c *Creature) *IncreaseDefenseModifier {
	return &IncreaseDefenseModifier{CreatureModifier{creature: c}}
}

func (i *IncreaseDefenseModifier) Apply() {
	if i.creature.Attack <= 200 {
		fmt.Println("Increasing ", i.creature.Attack, " defense")
		i.creature.Defense++
	}

	i.CreatureModifier.Apply()
}

/*
	No bonuses modifier - Apply() forgets to call next.Apply() that means every single modifier after this is not going to be applied
*/
type NoBonusesModifier struct {
	CreatureModifier
}

func NewNoBonusesModifier(c *Creature) *NoBonusesModifier {
	return &NoBonusesModifier{CreatureModifier{creature: c}}
}

func (n *NoBonusesModifier) Apply() {}

func main() {
	goblin := NewCreature("Goblin", 1, 1)
	rootNode := NewCreatureModifier(goblin)       // first modifier
	rootNode.Add(NewDoubleAttackModifier(goblin)) // second modifier
	rootNode.Apply()                              // applies all the modifiers from LL
	fmt.Println(goblin.String())
	rootNode.Add(NewDoubleAttackModifier(goblin))
	rootNode.Add(NewDoubleAttackModifier(goblin))
	rootNode.Apply()
	fmt.Println(goblin.String())

	rootNode.Add(NewIncreasedDefenseModifier(goblin)) // runs all the modifiers on rootNode again, on previous value - root first modifier -----> latest modifier
	rootNode.Apply()
	fmt.Println(goblin.String())

	goblin2 := NewCreature("Goblin", 1, 1)
	rootNode2 := NewCreatureModifier(goblin2)           // first modifier
	rootNode2.Add(NewDoubleAttackModifier(goblin2))     // second modifier
	rootNode2.Add(NewNoBonusesModifier(goblin2))        // no bonus modifier - aage ke sare modifiers will be disabled
	rootNode2.Add(NewIncreasedDefenseModifier(goblin2)) // fourth modifier
	rootNode2.Apply()
	fmt.Println(goblin2.String())
}
