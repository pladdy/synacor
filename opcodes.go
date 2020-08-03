package main

import (
	"fmt"
	"os"
)

type opcode uint8

const (
	opHalt opcode = iota // 0
	opSet                // 1
	opPush               // 2
	opPop                // 3
	opEq                 // 4
	opGt                 // 5
	opJmp                // 6
	opJt                 // 7
	opJf                 // 8
	opAdd                // 9
	opMult               // 10
	opMod                // 11
	opAnd                // 12
	opOr                 // 13
	opNot                // 14
	opRmem               // 15
	opWmem               // 16
	opCall               // 17
	opRet                // 18
	opOut                // 19
	opIn                 // 20
	opNoop               // 21
)

type operator func(i int, m *[]uint16, r *registers) int

var operatorMap = map[opcode]operator{
	opHalt: halt,
	opSet:  set,
	opPush: push,
	opPop:  notImplemented,
	opEq:   eq,
	opGt:   notImplemented,
	opJmp:  jump,
	opJt:   jumpTrue,
	opJf:   jumpFalse,
	opAdd:  add,
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

// template
// func bloop(i int, m *[]uint16, r *registers) int {
// 	return i
// }

// add: 9 a b c
//  assign into <a> the sum of <b> and <c> (modulo 32768)
func add(i int, m *[]uint16, r *registers) int {
	i, a := getNextRawValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)
	i, c := getNextValueShiftIndex(i, m, r)
	r.Set(a, b+c)
	return i
}

// eq: 4 a b c
//   set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func eq(i int, m *[]uint16, r *registers) int {
	i, a := getNextRawValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)
	i, c := getNextValueShiftIndex(i, m, r)

	if b == c {
		r.Set(a, 1)
	} else {
		r.Set(a, 0)
	}

	return i
}

// Similar to getNextValueShiftIndex, but returns the raw value, not the value
// in the register if the value is a register
func getNextRawValueShiftIndex(i int, m *[]uint16, r *registers) (int, uint16) {
	newIndex := i + 1
	return newIndex, (*m)[newIndex]
}

// This returns the value and shifts the provided index...
// TODO: This function is doing too much...what if memory was a struct that
//       handled the index value?
func getNextValueShiftIndex(i int, m *[]uint16, r *registers) (int, uint16) {
	newIndex := i + 1
	value := (*m)[newIndex]

	// what if the index is a register we want to set a value to?
	if isRegister(value) {
		value = r.Get(value)
	}
	return newIndex, value
}

// halt: 0
//   stop execution and terminate the program
func halt(i int, m *[]uint16, r *registers) int {
	os.Exit(0)
	return i
}

// jmp: 6 a
//   jump to <a>
func jump(i int, m *[]uint16, r *registers) int {
	_, jumpLocation := getNextValueShiftIndex(i, m, r)
	// After the jump is called the calling loop will iterate and increment
	// the index.  This decrements the index in preperation for that.
	// TODO: is there a better way to do this?
	//   - IE should we just call operatorMap[int(jumpLocation)]...?
	return int(jumpLocation) - 1
}

// jf: 8 a b
//   if <a> is zero, jump to <b>
func jumpFalse(i int, m *[]uint16, r *registers) int {
	i, a := getNextValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)

	if a == 0 {
		return int(b) - 1
	}
	return i
}

// jt: 7 a b
//   if <a> is nonzero, jump to <b>
func jumpTrue(i int, m *[]uint16, r *registers) int {
	i, a := getNextValueShiftIndex(i, m, r)
	i, b := getNextValueShiftIndex(i, m, r)

	if a > 0 {
		return int(b) - 1
	}
	return i
}

// noop: 21
//   no operation
func noop(i int, m *[]uint16, r *registers) int {
	return i
}

func notImplemented(i int, m *[]uint16, r *registers) int {
	panic(fmt.Sprintf("opCode %d not implemented", (*m)[i]))
}

// set: 1 a b
//   set register <a> to the value of <b>
func set(i int, m *[]uint16, r *registers) int {
	i, a := getNextRawValueShiftIndex(i, m, r)
	i, b := getNextRawValueShiftIndex(i, m, r)

	// assume it's a register...
	r.Set(a, b)
	return i
}

// out: 19 a
//   write the character represented by ascii code <a> to the terminal
func out(i int, m *[]uint16, r *registers) int {
	index, value := getNextValueShiftIndex(i, m, r)
	fmt.Print(string(value))
	return index
}

// push: 2 a
//   push <a> onto the stack
func push(i int, m *[]uint16, r *registers) int {
	return i
}
