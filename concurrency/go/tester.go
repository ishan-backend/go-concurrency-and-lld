package main

import (
	"./patterns/cancellation"
	"runtime"
)

func main() {
	// This function controls the number of operating system threads that will host so-called “Work Queues.”
	runtime.GOMAXPROCS(runtime.NumCPU())

	//_buffered_channel.TestChannelOwner()
	//_select.TestSelect()

	/*
		Concurrency patterns
	*/
	//_confinement.TestAdhoc()
	//_confinement.TestLexical1()

	_cancellation.TestCancellation()

}
