package _confinement

import (
	"bytes"
	"fmt"
	"sync"
)

// https://github.com/ishan-backend/go-concurrency-guide#confinement
// ensuring information is only ever available from one concurrent process

// TestAdhoc achieves confinement through a convention
func TestAdhoc() {
	data := make([]int, 4)
	data = []int{1, 2, 3, 4}

	informationCh := make(chan int)
	go func() {
		defer close(informationCh)
		for i := range data {
			informationCh <- data[i]
		}
	}()

	for d := range informationCh {
		fmt.Printf("data entity %v \n", d)
	}
}

// TestLexical1 : Lexical confinement involves using lexical scope to expose only the correct data and concurrency primitives for multiple concurrent processes to use.
// one to many consumers
func TestLexical1() {
	data := []byte("golang")
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buffer bytes.Buffer
		for _, c := range data {
			fmt.Fprintf(&buffer, "%c", c)
		}
		fmt.Println(buffer.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}

// TestLexical2 one to many consumers
func TestLexical2() {
	// producer
	_ = func() <-chan int {
		result := make(chan int, 5)
		go func() {
			defer close(result)

			for i := 0; i < 5; i++ {
				result <- i
			}
		}()

		return result
	}
}
