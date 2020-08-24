package main

import "github.com/pladdy/synacor"

func main() {
	m := synacor.NewMachine()
	m.Load("./challenge.bin")
	m.Run()
}
