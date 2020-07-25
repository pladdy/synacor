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

type operator func(i int, m *[]uint16) int

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

func halt(i int, m *[]uint16) int {
	os.Exit(0)
	return i
}

func noop(i int, m *[]uint16) int {
	return i
}

func notImplemented(i int, m *[]uint16) int {
	fmt.Println("opCode not implemented:", (*m)[i])
	panic("not implemented")
}

func out(i int, m *[]uint16) int {
	arg1Index := i + 1
	fmt.Print(string((*m)[arg1Index]))
	return arg1Index
}
