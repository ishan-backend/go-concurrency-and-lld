package _cancellation

import (
	"fmt"
	"sync"
	"time"
)

func TestCancellation() {

	doWork := func(doneCh <-chan interface{}, data <-chan string, wg *sync.WaitGroup) chan interface{} {
		terminatedCh := make(chan interface{})
		wg.Add(1)

		go func() {
			defer fmt.Println("... cancelled/terminated doWork go routine")
			defer close(terminatedCh)
			defer wg.Done()

			for {
				select {
				case <-doneCh:
					fmt.Println("closing via done channel val")
					return // on return, control flow will break the while loop, go routine statements will be finished so defer statements will be called

				/*
					case <-data:
						- will not cause deadlock only if nil is passed as argument, because it won't wait for channel value and cause deadlock

					in case when you have to pass some value in channel, pass it via another go-routine so that this go-routine is not waiting in deadlock

				*/
				case s, ok := <-data:
					if !ok {
						fmt.Println("irregular value wait: ", s, ok)
						continue
					}
					fmt.Println("data is: ", s)
				}
			}
		}()

		return terminatedCh
	}

	doneCh := make(chan interface{})
	dataCh := make(chan string) // unbuffered data channel - This means that an unbuffered channel can only hold one value at a time. When a sender sends a value on an unbuffered channel, it will block until there is a receiver ready to receive the value.

	var wg sync.WaitGroup
	_ = doWork(doneCh, dataCh, &wg) // reference to the channel is returned implicitly

	go func() {
		time.Sleep(1 * time.Second)
		dataCh <- "ping" // Send data to the channel
		close(dataCh)    // Close the channel after sending data
	}()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("cancelling doWork go-routine ...")
		close(doneCh) // as soon as doneCh is closed, a signal goes to channel, and select statement from other go-routine gets ready
	}()

	wg.Wait()
	fmt.Println("closing entire application !!! .....")
}
