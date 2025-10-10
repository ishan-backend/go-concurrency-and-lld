package main

import (
	"fmt"
	multithreadedParser "lld/multithreaded_paranthesis_parser"
)

func main() {
	fmt.Println(multithreadedParser.ValidateMultithreaded("(()()())("))
}
