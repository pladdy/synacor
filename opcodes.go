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
	opJt:   jumpTrue,
	opJf:   jumpFalse,
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
func getNextValueShiftIndex(i int, m *[]uint16, r *registers) (int, uint16) {
	newIndex := i + 1
	value := (*m)[newIndex]
	//fmt.Println("  getNextValue called, newIndex:", newIndex, "value:", value)

	if isRegister(value) {
		value = r.Get(value)
	}
	return newIndex, value
}

func halt(i int, m *[]uint16, r *registers) int {
	os.Exit(0)
	return i
}

func jump(i int, m *[]uint16, r *registers) int {
	_, jumpLocation := getNextValueShiftIndex(i, m, r)
	// fmt.Println("  Jump has been called, location is:", jumpLocation)

	// After the jump is called the calling loop will iterate and increment
	// the index.  This decrements the index in preperation for that.
	// TODO: is there a better way to do this?
	//   - IE should we just call operatorMap[int(jumpLocation)]...?
	return int(jumpLocation) - 1
}

func jumpFalse(i int, m *[]uint16, r *registers) int {
	i, a := getNextValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)

	if a == 0 {
		// return jump(int(b), m, r)
		return int(b) - 1
	}
	return i
}

func jumpTrue(i int, m *[]uint16, r *registers) int {
	i, a := getNextValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)

	if a > 0 {
		// return jump(int(b), m, r)
		return int(b) - 1
	}
	return i
}

func noop(i int, m *[]uint16, r *registers) int {
	return i
}

func notImplemented(i int, m *[]uint16, r *registers) int {
	panic(fmt.Sprintf("opCode %d not implemented", (*m)[i]))
}

func out(i int, m *[]uint16, r *registers) int {
	index, value := getNextValueShiftIndex(i, m, r)
	fmt.Print(string(value))
	return index
}
