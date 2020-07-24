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

type operator func(o opcode, i int, m memoryStack) int

var operatorMap = map[opcode]operator{
	opHalt: halt,
	opSet:  notImplemented,
	opPush: notImplemented,
	opPop:  notImplemented,
	opEq:   notImplemented,
	opGt:   notImplemented,
	opJmp:  notImplemented,
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

func halt(o opcode, i int, m memoryStack) int {
	os.Exit(0)
	return i
}

func noop(o opcode, i int, m memoryStack) int {
	return i
}

func notImplemented(o opcode, i int, m memoryStack) int {
	panic("not implemented")
}

func out(o opcode, i int, m memoryStack) int {
	v := m[i]
	fmt.Print(string(v))
	return i + 1
}
