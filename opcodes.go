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

type operator func(p *program, r *registers)

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
// func bloop(p *program, r *registers) {
// 	// }

// add: 9 a b c
//  assign into <a> the sum of <b> and <c> (podulo 32768)
func add(p *program, r *registers) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.Set(a, b+c)
}

// eq: 4 a b c
//   set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func eq(p *program, r *registers) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)

	if b == c {
		r.Set(a, 1)
	} else {
		r.Set(a, 0)
	}
}

// halt: 0
//   stop execution and terminate the program
func halt(p *program, r *registers) {
	os.Exit(0)
}

// jmp: 6 a
//   jump to <a>
func jump(p *program, r *registers) {
	p.index = int(p.getNext(r))
}

// jf: 8 a b
//   if <a> is zero, jump to <b>
func jumpFalse(p *program, r *registers) {
	a := p.getNext(r)
	b := p.getNext(r)

	if a == 0 {
		p.index = int(b)
	} else {
		p.index = p.index + 1
	}
}

// jt: 7 a b
//   if <a> is nonzero, jump to <b>
func jumpTrue(p *program, r *registers) {
	a := p.getNext(r)
	b := p.getNext(r)

	if a > 0 {
		p.index = int(b)
	} else {
		p.index = p.index + 1
	}
}

// noop: 21
//   no operation
func noop(p *program, r *registers) {
	p.index = p.index + 1
}

func notImplemented(p *program, r *registers) {
	panic(fmt.Sprintf("opCode %d not implemented", p.memory[p.index]))
}

// out: 19 a
//   write the character represented by ascii code <a> to the terminal
func out(p *program, r *registers) {
	fmt.Print(string(p.getNext(r)))
	p.index = p.index + 1
}

// push: 2 a
//   push <a> onto the stack
func push(p *program, r *registers) {
}

// set: 1 a b
//   set register <a> to the value of <b>
func set(p *program, r *registers) {
	a := p.getNextRaw()
	b := p.getNextRaw()
	// assume it's a register...
	r.Set(a, b)
}
