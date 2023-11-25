package _buffered_channel

import (
	"fmt"
	"sync"
)

/*
Table with result of channel operations

Operation	Channel State	Result

Read	nil	Block
_	Open and Not Empty	Value
_	Open and Empty	Block
_	Close	default value, false
_	Write Only	Compilation Error
Write	nil	Block
_	Open and Full	Block
_	Open and Not Full	Write Value
_	Closed	panic
_	Receive Only	Compilation Error
Close	nil	panic
_	Open and Not Empty	Closes Channel; reads succeed until channel is drained, then reads produce default value
_	Open and Empty	Closes Channel; reads produces default value
_	Closed	panic

TIP: Cannot close a receive-only channel
Here, Block = deadlocking
*/

type ChannelOwner struct {
	writeCh   chan int      // nil channel
	readerCh  chan struct{} // nil channel to signal completion, using bidirectional as Receive-only channel cannot be closed
	closeOnce sync.Once     // ensures close() is called only once
}

func NewChannelOwner() *ChannelOwner {
	return &ChannelOwner{
		writeCh:  make(chan int),      // open and not full channel; bidirectional (read via range; write allowed)
		readerCh: make(chan struct{}), // open and empty channel; unidirectional
	}
}

func (c *ChannelOwner) Start() {
	go func() {
		// Perform Writes
		// There can also be case where you can pass ownership of writing to this writeCh to another go-routine
		for i := 0; i < 5; i++ {
			c.writeCh <- i
		}

		// Close write channel
		c.closeOnce.Do(func() {
			close(c.writeCh)
		})
	}()
}

func (c *ChannelOwner) TransferOwnership(newCh *chan int, wg *sync.WaitGroup) {
	// pass old channel reference (maybe closed already) to new channel reference, owned by different goroutine
	// cannot write to closed channel reference
	c.closeOnce.Do(func() {
		close(c.writeCh)
		*newCh = c.writeCh
		c.writeCh = nil
	})

	close(c.readerCh)
	wg.Done()
}

func TestChannelOwner() {
	channelOwner := NewChannelOwner()
	channelOwner.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		newCh := make(chan int)
		channelOwner.TransferOwnership(&newCh, &wg)
		// read from transferred channel
		for v := range newCh {
			fmt.Println("Received in transferred goroutine", v)
		}
	}()

	wg.Wait()
	// This is often used as a synchronization mechanism. The <-c.readerCh line in your TestChannelOwner function will unblock when the c.readerCh channel is closed, allowing the program to proceed.
	<-channelOwner.readerCh
}
