package main

// Scalar object
type Neuron struct {
	In, Out []*Neuron
}

func (n *Neuron) ConnectTo(other *Neuron) {
	n.Out = append(n.Out, other)
	other.In = append(other.In, n)
}

// Composite Object
type NeuronLayer struct {
	Neurons []Neuron
}

func NewNeuronLayer(count int) *NeuronLayer { // factory for NeuronLayer
	return &NeuronLayer{Neurons: make([]Neuron, count)}
}

type NeuronInterface interface {
	Iter() []*Neuron // gives us pointers to all the neurons in particular type of object (Neuron, NeuronLayer)
}

func (n *Neuron) Iter() []*Neuron {
	res := make([]*Neuron, 0)
	res = append(res, n)
	return res
}

func (n *NeuronLayer) Iter() []*Neuron {
	res := make([]*Neuron, 0)
	for i := range n.Neurons {
		res = append(res, &n.Neurons[i])
	}
	return res
}

// Connect - NeuronInterface helps scalar object masquerade as composite object
func Connect(left, right NeuronInterface) {
	for _, l := range left.Iter() {
		for _, r := range right.Iter() {
			l.ConnectTo(r)
		}
	}
}

func main() {
	neuron1, neuron2 := &Neuron{}, &Neuron{}
	layer1, layer2 := NewNeuronLayer(3), NewNeuronLayer(4)

	Connect(neuron1, neuron2)
	Connect(neuron1, layer1)
	Connect(neuron2, layer2)
	Connect(neuron2, layer1)
}
