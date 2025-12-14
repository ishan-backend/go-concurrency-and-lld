package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

/*
Critical Point: When main() returns, the entire program exits, even if other goroutines are still running.
Channels : https://www.concurrency.rocks/tutorials/channels

// Consider buffer size carefully

	func buffered() <-chan int {
	    ch := make(chan int, 100)  // Large buffer
	    // Pros: Less blocking
	    // Cons: More memory, delayed error detection
	}

Small buffers (1-10) often sufficient
Benchmark buffer sizes for your use case
Consider batching for high-frequency generators
Use sync.Pool for temporary objects in generators
*/
func CloseChannels() {
	// closing channels

	/*
		Only the sender should close a channel
		Sending on a closed channel causes panic
		Receiving from a closed channel returns zero value
		Use v, ok := <-ch to check if channel is closed OR You can use range to receive values until a channel is closed:
	*/
	jobs := make(chan int, 5) // buffered channel, sends are not blocked if the channel is full, receives are blocked if the channel is empty
	done := make(chan bool)

	go func() {
		if true {
			for {
				j, more := <-jobs
				if !more { // if the channel is closed, the more will be false
					fmt.Println("received all jobs")
					done <- true
					return
				}
				fmt.Println("received job", j)
			}
		} else {
			for elem := range jobs {
				fmt.Println("received job", elem)
			}
			fmt.Println("received all jobs")
			done <- true
			return
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
		time.Sleep(time.Second)
	}
	close(jobs) // closing the channel, no more sends are allowed, receives are not blocked
	fmt.Println("sent all jobs")

	<-done
}

// GeneratorPattern are functions that return channels
func GeneratorPattern() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i // unbuffered channel, ingestion is blocked until is received, hence it is taken out from channel
		}
		close(ch)
	}()

	return ch
}

func GeneratorButConsumerControlsTermination(done <-chan bool) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)

		/*
			No priority ordering: select doesn't prioritize cases by order.
			If multiple cases are ready, one is chosen pseudo-randomly (to prevent starvation).
			If only one case is ready, that one is chosen.
			If none are ready, select blocks until at least one becomes ready.
		*/
		for i := 0; ; i++ {
			select {
			case ch <- i: // ch is unbuffered, this will execute only when there is a receiver
			case <-done: // in case ch is blocked, this condition is tested
				return
			}
		}

	}()

	return ch
}

func HelperGeneratorButConsumerControlsTermination() {
	done := make(chan bool)
	ch := GeneratorButConsumerControlsTermination(done)

	// Consumer receives a few values
	for i := 0; i < 3; i++ {
		val := <-ch
		fmt.Printf("Received: %d\n", val)
		time.Sleep(100 * time.Millisecond)
	}

	// Signal done - this will cause the generator to exit
	close(done) // or: done <- true

	// Try to receive remaining values (if any)
	// Since generator exits when done is closed, channel will be closed
	for val := range ch {
		fmt.Printf("Received after done: %d\n", val)
	}

	fmt.Println("Generator terminated")
}

// Combine multiple generators into one
func FanInGenerator(input1, input2 <-chan int) <-chan int {
	ch := make(chan int)

	/*
		Uses select instead of sequential receives: receives from whichever channel has data, avoiding blocking.
		Tracks closed channels: uses local variables in1 and in2 that can be set to nil when closed.
		Nil channel idiom: when a channel closes, set it to nil so select ignores it (Go idiom).
		Proper termination: checks if both channels are nil (closed) and exits the goroutine, closing the output channel via defer.

	*/
	go func() {

		// Use local variables so we can set them to nil when closed
		in1 := input1
		in2 := input2

		for {
			// If both channels are nil (closed), exit
			if in1 == nil && in2 == nil {
				close(ch)
				return
			}

			select {
			case s, ok := <-in1:
				if !ok {
					// Channel is closed, set to nil so select ignores it
					in1 = nil
					continue
				}
				ch <- s

			case s, ok := <-in2:
				if !ok {
					// Channel is closed, set to nil so select ignores it
					in2 = nil
					continue
				}
				ch <- s
			}
		}
	}()

	return ch
}

func RateLimitedGenerator() {

}

type Event struct {
	Type string
	Time time.Time
}

func EventStreamGenerator(done <-chan bool) <-chan Event {
	ch := make(chan Event, 10) // keeping small buffers allow non-blocking event stream; when done arrives then production is paused
	go func() {
		for {
			select {
			case <-done:
				close(ch)
				return
			default:
				ch <- Event{
					Type: "click",
					Time: time.Now(),
				}
			}
		}
	}()

	return ch
}

func FileLineGenerator(path string) (<-chan string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch, nil
}

// context timeout implementation - https://www.concurrency.rocks/tutorials/first-racer
func withTimeout(fn func() string, timeout time.Duration) (string, error) {
	ch := make(chan string, 1)

	go func() {
		ch <- fn()
	}()

	select {
	case result := <-ch:
		return result, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout after %v", timeout)
	}
}

// ------------------------------------------------------------

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []string
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Put(item string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	q.cond.Signal() // Wake one waiter
}

func (q *Queue) Get() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait() // Releases mutex while waiting
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// ------------------------------------------------------------

// Better: Use sync.Once
// type BetterSingleton struct {
//     once sync.Once
//     data *expensiveResource
// }

// func (s *BetterSingleton) Get() *expensiveResource {
//     s.once.Do(func() {
//         s.data = createExpensiveResource()
//     })
//     return s.data
// }

// WaitGroup approach
func withWaitGroup(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}
	wg.Wait()
}

// Channel approach
func withChannels(tasks []func()) {
	done := make(chan bool, len(tasks))
	for _, task := range tasks {
		go func(t func()) {
			t()
			done <- true
		}(task)
	}
	for range tasks {
		<-done
	}
}

func main() {
	/*
		CloseChannels()
		fmt.Println("Tested Close Channels")

		ch := GeneratorPattern()
		for el := range ch {
			fmt.Printf("generated %v", el)
		}
		fmt.Println("generator pattern - return channel tested")
	*/

	// HelperGeneratorButConsumerControlsTermination()

	// ch := FanInGenerator(GeneratorPattern(), GeneratorPattern())
	// for el := range ch { // only pushed values are received
	// 	fmt.Printf("%v", el)
	// }
	// fmt.Println("closed fan in generator")

	// done := make(chan bool)
	// ch := EventStreamGenerator(done)
	// time.Sleep(4000 * time.Duration(time.Millisecond))
	// close(done)
	// for el := range ch {
	// 	fmt.Printf("%v", el)
	// }

}
