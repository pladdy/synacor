package main

import (
	"fmt"
	"os"
)

type opcode uint8

const (
	opHalt opcode = iota
	opSet
	opPush
	opPop
	opEq
	opGt
	opJmp
	opJt
	opJf
	opAdd
	opMult
	opMod
	opAnd
	opOr
	opNot
	opRmem
	opWmem
	opCall
	opRet
	opOut
	opIn
	opNoop
)

func noop() {
	return
}

func halt() {
	os.Exit(0)
}

func out(c uint16) {
	fmt.Print(string(c))
}
