package main

import (
	"fmt"
	"strings"
)

const (
	indentSize = 2
)

type HTMLElement struct {
	name, text string
	elements   []HTMLElement
}

func (e *HTMLElement) String() string {
	return e.string(0) // pass indents required for one HTML element say <li> </li>
}

func (e *HTMLElement) string(indent int) string { // recursive function
	sb := strings.Builder{}
	i := strings.Repeat(" ", indent*indentSize)        // build total indent ki empty string
	sb.WriteString(fmt.Sprintf("%s<%s>\n", i, e.name)) // <li>
	if len(e.text) > 0 {
		sb.WriteString(strings.Repeat(" ", (indent+1)*indentSize)) // add additional indent space before writing text
		sb.WriteString(e.text)
		sb.WriteString("\n")
	}

	// for all internal elements
	for _, el := range e.elements {
		sb.WriteString(el.string(indent + 1))
	}

	sb.WriteString(fmt.Sprintf("%s</%s>\n", i, e.name)) // </li>
	return sb.String()
}

type HTMLBuilder struct {
	root     HTMLElement // you can call root.String() to print HTML in console
	rootName string      // sometimes you would want to clear the builder, so we need this in cache as an identifier
}

func NewHTMLBuilder(rootName string) *HTMLBuilder {
	return &HTMLBuilder{
		root:     HTMLElement{rootName, "", []HTMLElement{}},
		rootName: rootName,
	}
}

func (b *HTMLBuilder) String() string { // prints all the HTML elements from builder root and children if any
	return b.root.String()
}

func (b *HTMLBuilder) AddChild(childName, childText string) {
	e := HTMLElement{
		name:     childName,
		text:     childText,
		elements: []HTMLElement{},
	}
	b.root.elements = append(b.root.elements, e)
}

func (b *HTMLBuilder) AddChildFluent(childName, childText string) *HTMLBuilder {
	e := HTMLElement{
		name:     childName,
		text:     childText,
		elements: []HTMLElement{},
	}
	b.root.elements = append(b.root.elements, e)
	return b // returns pointer to receiver helps user to reuse results of a call to continue calling on the object
}

func main() {
	// end user only needs to care about utility functions
	b := NewHTMLBuilder("ul")
	b.AddChild("li", "hello")
	b.AddChild("li", "world")
	fmt.Println(b.String())

	// chain calls in builder pattern
	b = NewHTMLBuilder("ul")
	b.AddChildFluent("li", "hello1").AddChildFluent("li", "world1")
	fmt.Println(b.String())
}
