package main

import (
	"fmt"
	"strings"
)

type OutputFormat int

const ( // Concrete implementation to make
	Markdown OutputFormat = iota
	HTML
)

/*
	HTML:
		<ul>
			<li> </li>
			<li> </li>
		</ul>

	Markdown:
	*
		1.
		2.
	*
*/

type ListStrategy interface { // interface which we expect each strategy to follow
	Start(builder *strings.Builder)
	End(builder *strings.Builder)
	AddItem(builder *strings.Builder, item string)
}

type MarkdownListStrategy struct{}

func (m *MarkdownListStrategy) Start(builder *strings.Builder) {
}

func (m *MarkdownListStrategy) End(builder *strings.Builder) {
}

func (m *MarkdownListStrategy) AddItem(builder *strings.Builder, item string) {
	builder.WriteString("* " + item + "\n")
}

type HtmlListStrategy struct{}

func (m *HtmlListStrategy) Start(builder *strings.Builder) {
	builder.WriteString("<ul>\n")
}

func (m *HtmlListStrategy) End(builder *strings.Builder) {
	builder.WriteString("</ul>\n")
}

func (m *HtmlListStrategy) AddItem(builder *strings.Builder, item string) {
	builder.WriteString("  <li>" + item + "</li>\n")
}

type TextProcessor struct { // support injection of strategy in high level algorithm
	builder  strings.Builder
	strategy ListStrategy // ListStrategy is implemented by concrete implementations
}

func NewTextProcessor(strategy ListStrategy) *TextProcessor {
	return &TextProcessor{
		builder:  strings.Builder{},
		strategy: strategy,
	}
}

func (t *TextProcessor) SetOutputFormat(of OutputFormat) {
	switch of {
	case Markdown:
		t.strategy = &MarkdownListStrategy{}
	case HTML:
		t.strategy = &HtmlListStrategy{}
	}
}

func (t *TextProcessor) AppendItemsToList(items []string) {
	s := t.strategy
	s.Start(&t.builder)
	for _, i := range items {
		s.AddItem(&t.builder, i)
	}
	s.End(&t.builder)
}

func (t *TextProcessor) Reset() {
	t.builder.Reset()
}

func (t *TextProcessor) String() string {
	return t.builder.String()
}

func main() {
	tp := NewTextProcessor(&MarkdownListStrategy{})
	tp.AppendItemsToList([]string{"foo", "bar"})
	fmt.Println(tp.String())

	tp.Reset()
	tp.SetOutputFormat(HTML)
	tp.AppendItemsToList([]string{"foo", "bar"})
	fmt.Println(tp.String())
}
