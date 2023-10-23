package main

import "fmt"

type Lizard struct {
	Age int
}

func (l *Lizard) Crawl() {
	fmt.Println("Lizard is crawling")
}

type Bird struct {
	Age int
}

func (b *Bird) Fly() {
	if b.Age >= 10 {
		fmt.Println("Bird is flying")
	}
}

type Dragon struct {
	Bird
	Lizard
}

func (d *Dragon) setAge(age int) {
	d.Bird.Age = age
	d.Lizard.Age = age
}

func (d *Dragon) getAge() int {
	return d.Bird.Age
}

func main() {
	d := Dragon{}
	//d.Age = 10
	d.Bird.Age = 10
	d.Lizard.Age = 10
	d.Fly()
	d.Crawl()

	// you need to manage the age state consistent for a dragon, even with getter and setter for age
	d2 := Dragon{}
	d2.setAge(10)
	d2.Fly()
	d2.Crawl()
	d2.Bird.Age = 5
	d2.Fly() // ignored since age state is changed

}
