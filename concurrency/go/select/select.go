package _select

import (
	"fmt"
	"time"
)

/*
	the select cases do not work the same as the switch, which is sequential, and the execution will not automatically fall if none of the criteria is met.
	instead, all channel reads and writes are considered simultaneously to see if any of them are ready: channels filled or closed in the case of reads and channels not at capacity in the case of writes.
    if none of the channels are ready, the entire select command is blocked. Then, when one of the channels is ready, the operation will proceed and its corresponding instructions will be executed.

	When using select, it's important to understand that the cases are considered in the order they appear. In your code, the case <-time.After(2 * time.Second) is the last case, and it will only be considered if none of the other cases are ready.
*/

func TestSelect() {
	currTime := time.Now()
	ch := make(chan interface{})
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			ch <- i
		}

		close(ch)
	}()

	timeout := time.After(15 * time.Second) // poll multiple times on same timer in select

	fmt.Println("Blocking on read (select using for with channel) at ", currTime)
	// loop over the select statement to continuously receive values, otherwise select will be blocked till it receives only the first value in first case
loop:
	for {
		select {
		/*
				_, ok := <-time.After(10 * time.Second): - is called in default case, which means it will be selected only when none other cases are ready
				in our case, since ch is continuously sending values, time.After case is not selected until ch channel is closed

				to make timeout work as expected, we need to move time.After case to the same level as the case val, ok := <-ch:
				43

			    - Each time you execute time.After(4 * time.Second) you create a new timer channel. There's no way the select statement can remember the channel it selected on in the previous iteration

				case <-time.After(3 * time.Second):
					// Timed out handling - cannot complete all the work in 20 seconds and exiting the for loop
					fmt.Println("timed out ...")
					break loop
		*/
		case <-timeout:
			fmt.Println("timed out ...")
			break loop

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
			time.Sleep(1 * time.Second)
			fmt.Println("some non-blocking work")
		}
	}

furtherWork:
	fmt.Printf("Select statement exited successfully at %v \n", time.Now())

	/*
		To block forever:
		for {
			select {}
		}
	*/
}
