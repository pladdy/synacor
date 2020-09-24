package main

import (
	"fmt"

	"github.com/pladdy/synacor"
)

func main() {
	m := synacor.NewMachine()
	m.Load("./challenge.bin")

	i := 0
	for m.HasMoreOps() {
		opname, opcode, args := m.NextOp()
		fmt.Printf("Index: %d, Operation: %s (%d), Args: %d", i, opname, opcode, args)
		if opcode == 19 {
			fmt.Printf(" %s", string(args[0]))
		}
		fmt.Println()
		i = i + 1
	}
}
