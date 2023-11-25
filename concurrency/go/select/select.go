package _select

import (
	"fmt"
	"time"
)

/*
	the select cases do not work the same as the switch, which is sequential, and the execution will not automatically fall if none of the criteria is met.
	instead, all channel reads and writes are considered simultaneously to see if any of them are ready: channels filled or closed in the case of reads and channels not at capacity in the case of writes.
    if none of the channels are ready, the entire select command is blocked. Then, when one of the channels is ready, the operation will proceed and its corresponding instructions will be executed.
*/

func TestSelect() {
	currTime := time.Now()
	ch := make(chan interface{})

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			ch <- i
		}

		close(ch)
	}()

	fmt.Println("Blocking on read (select using for with channel) at ", currTime)
	// loop over the select statement to continuously receive values
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				// channel is closed, exit the while
				fmt.Println("channel closed received")
				// using return here, exits the program, if you want to continue work beyond breaking the while loop use goto
				goto furtherWork
			}
			fmt.Printf("Unblocked after %v , value received: %v \n", time.Since(currTime), val)

		default:
			// do other non-blocking work here
		}
	}

furtherWork:
	fmt.Printf("Select statement exited successfully at %v \n", time.Now())
}
