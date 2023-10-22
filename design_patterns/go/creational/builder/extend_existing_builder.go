package main

import "fmt"

type human struct {
	name, email string
	// newly added var
	positionName string
}

type humanMod func(*human)
type HumanBuilder struct {
	actions []humanMod
}

func (h *HumanBuilder) Called(name string) *HumanBuilder {
	// instead of just performing modification, we add a function to the list of actions, which will be called at runtime (FP)
	h.actions = append(h.actions, func(p *human) {
		p.name = name
	})
	return h
}

func (h *HumanBuilder) Build() *human {
	p := human{}
	for _, ac := range h.actions {
		ac(&p)
	}
	return &p
}

// extensibility for other domains - without using aggregation of different builders using FP approach
func (h *HumanBuilder) WorksAs(positionName string) *HumanBuilder {
	h.actions = append(h.actions, func(p *human) {
		p.positionName = positionName
	})
	return h
}

func main() {
	b := &HumanBuilder{}
	p := b.Called("Dmitri").WorksAs("CEO").Build()
	fmt.Println(p)
}
