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

type operator func(p *program, r *registers, s *stack)

var operatorMap = map[opcode]operator{
	opHalt: halt,
	opSet:  set,
	opPush: push,
	opPop:  pop,
	opEq:   eq,
	opGt:   gt,
	opJmp:  jump,
	opJt:   jumpTrue,
	opJf:   jumpFalse,
	opAdd:  add,
	opMult: notImplemented,
	opMod:  notImplemented,
	opAnd:  and,
	opOr:   or,
	opNot:  not,
	opRmem: notImplemented,
	opWmem: notImplemented,
	opCall: notImplemented,
	opRet:  notImplemented,
	opOut:  out,
	opIn:   notImplemented,
	opNoop: noop,
}

// add: 9 a b c
//  assign into <a> the sum of <b> and <c> (podulo 32768)
func add(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b+c)
	p.index = p.index + 1
}

// and: 12 a b c
//   stores into <a> the bitwise and of <b> and <c>
func and(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b&c)
	p.index = p.index + 1
}

// eq: 4 a b c
//  .set <a> to 1 if <b> is equal to <c>;.set it to 0 otherwise
func eq(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)

	//fmt.Println("Eq, A:", a, "B:", b, "C:", c)

	if b == c {
		r.set(a, 1)
	} else {
		r.set(a, 0)
	}
	p.index = p.index + 1
}

// gt: 5 a b c
//   set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func gt(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)

	if b > c {
		r.set(a, 1)
	} else {
		r.set(a, 0)
	}
	p.index = p.index + 1
}

// halt: 0
//   stop execution and terminate the program
func halt(p *program, r *registers, s *stack) {
	os.Exit(0)
}

// jmp: 6 a
//   jump to <a>
func jump(p *program, r *registers, s *stack) {
	p.index = int(p.getNext(r))
}

// jf: 8 a b
//   if <a> is zero, jump to <b>
func jumpFalse(p *program, r *registers, s *stack) {
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
func jumpTrue(p *program, r *registers, s *stack) {
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
func noop(p *program, r *registers, s *stack) {
	p.index = p.index + 1
}

// not: 14 a b
//   stores 15-bit bitwise inverse of <b> in <a>
func not(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	r.set(a, ^b%modulo)
	p.index = p.index + 1
}

func notImplemented(p *program, r *registers, s *stack) {
	panic(fmt.Sprintf("opCode %d not implemented", p.memory[p.index]))
}

// or: 13 a b c
//   stores into <a> the bitwise or of <b> and <c>
func or(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b|c)
	p.index = p.index + 1
}

// out: 19 a
//   write the character represented by ascii code <a> to the terminal
func out(p *program, r *registers, s *stack) {
	fmt.Print(string(p.getNext(r)))
	p.index = p.index + 1
}

// push: 2 a
//   push <a> onto the stack
func push(p *program, r *registers, s *stack) {
	a := p.getNext(r)
	s.push(a)
	p.index = p.index + 1
}

// pop: 3 a
//   remove the top element from the stack and write it into <a>; empty stack = error
func pop(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := s.pop()
	r.set(a, b)
	p.index = p.index + 1
}

// set: 1 a b
//   set register <a> to the value of <b>
func set(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNextRaw()

	if isRegister(a) {
		r.set(a, b)
	}
	p.index = p.index + 1
}
