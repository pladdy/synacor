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

type operator func(i int, m *[]uint16, r *registers) int

var operatorMap = map[opcode]operator{
	opHalt: halt,
	opSet:  notImplemented,
	opPush: notImplemented,
	opPop:  notImplemented,
	opEq:   notImplemented,
	opGt:   notImplemented,
	opJmp:  jump,
	opJt:   notImplemented,
	opJf:   notImplemented,
	opAdd:  notImplemented,
	opMult: notImplemented,
	opMod:  notImplemented,
	opAnd:  notImplemented,
	opOr:   notImplemented,
	opNot:  notImplemented,
	opRmem: notImplemented,
	opWmem: notImplemented,
	opCall: notImplemented,
	opRet:  notImplemented,
	opOut:  out,
	opIn:   notImplemented,
	opNoop: noop,
}

// This returns the value and shifts the provided index...
// TODO: This function is doing too much
func getNextValueShiftIndex(i int, m *[]uint16) (int, uint16) {
	return i + 1, (*m)[i+1]
}

func halt(i int, m *[]uint16, r *registers) int {
	os.Exit(0)
	return i
}

func jump(i int, m *[]uint16, r *registers) int {
	_, jumpLocation := getNextValueShiftIndex(i, m)

	if isRegister(jumpLocation) {
		jumpLocation = r.Get(jumpLocation)
	}
	// After the jump is called the calling loop will iterate and increment
	// the index.  This decrements the index in preperation for that.
	// TODO: is there a better way to do this?
	//   - IE should we just call operatorMap[int(jumpLocation)]...?
	return int(jumpLocation) - 1
}

func noop(i int, m *[]uint16, r *registers) int {
	return i
}

func notImplemented(i int, m *[]uint16, r *registers) int {
	panic(fmt.Sprintf("opCode %d not implemented", (*m)[i]))
}

func out(i int, m *[]uint16, r *registers) int {
	index, value := getNextValueShiftIndex(i, m)
	fmt.Print(string(value))
	return index
}
