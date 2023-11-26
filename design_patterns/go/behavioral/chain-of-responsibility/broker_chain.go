package main

import "sync"

// CoR, Mediator, Observer, CQS

type Goku struct {
	// Mediator - central component which every Goku / Creature refers to; Since most Creatures participate in a game, its called game.
	game *Game

	Name            string
	attack, defense int // why lowercase? in method_chain, we applied Modifiers explicitly by calling Apply() method. Here we want to apply modifiers automatically, as soon as you make modifier
	// pass in Goku, then modifier automatically gets applied
	// so when you query Goku's attack and defense values, you get the final calculated values
	// attack, defense - only store initial values
	// we will add whole system of structs on top of this to make it possible
}

/*
 CQS - you want to get a Creature attack and defense values, you pass a Query struct to the Creature
 as opposed to calling method on it
*/

type Argument int

const (
	Attack Argument = iota
	Defense
)

type Query struct {
	CreatureName string
	WhatToQuery  Argument
	Value        int
}

/*
	Observer pattern:
*/
type Observer interface {
	Handle(q *Query) // query are events in this case
}

type Observable interface {
	Subscribe(o Observer)
	UnSubscribe(o Observer)
	Fire(q *Query) // when someone fires a query & that query gets processed by whoever is interested
}

/*
	Mediator pattern - game (centralised component)
*/
type Game struct {
	observers sync.Map // map of every single subscriber, and iterate through subscribers (map) and notify the subscriber that something has happened
}

func (g *Game) Subscribe(o Observer) {
	g.observers.Store(o, struct{}{})
}

func (g *Game) UnSubscribe(o Observer) {
	g.observers.Delete(o)
}

func (g *Game) Fire(q *Query) {
	g.observers.Range(func(key, value any) bool {
		if key == nil {
			return false
		}
		key.(Observer).Handle(q)
		return true
	})
}

/*
	Final implementation
*/
func NewGoku(game *Game, name string, attack, defense int) *Goku {
	return &Goku{
		game:    game,
		Name:    name,
		attack:  attack,
		defense: defense,
	}
}

// we do not address attack and defense directly, we have getters and setters

func main() {

}
