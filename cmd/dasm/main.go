package main

import (
	"fmt"

	"github.com/pladdy/synacor"
)

func main() {
	m := synacor.NewMachine()
	m.Load("./challenge.bin")
	fmt.Println(m.NextOp())
}
