package synacor

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

var operatorFunctionMap = map[opcode]operator{
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
	opMult: mult,
	opMod:  mod,
	opAnd:  and,
	opOr:   or,
	opNot:  not,
	opRmem: rmem,
	opWmem: wmem,
	opCall: call,
	opRet:  ret,
	opOut:  out,
	opIn:   in,
	opNoop: noop,
}

type operatorProperty struct {
	name string
	args int
}

var operatorPropertyMap = map[opcode]operatorProperty{
	opAdd:  {"add", 3},
	opAnd:  {"and", 3},
	opCall: {"call", 1},
	opEq:   {"eq", 3},
	opGt:   {"gt", 3},
	opHalt: {"halt", 0},
	opIn:   {"in", 1},
	opJmp:  {"jmp", 1},
	opJt:   {"jf", 2},
	opJf:   {"jt", 2},
	opMod:  {"mod", 3},
	opMult: {"mult", 3},
	opNoop: {"noop", 0},
	opNot:  {"not", 2},
	opOr:   {"or", 3},
	opOut:  {"out", 1},
	opPop:  {"pop", 1},
	opPush: {"push", 1},
	opRet:  {"ret", 0},
	opRmem: {"rmem", 2},
	opSet:  {"set", 2},
	opWmem: {"wmem", 2},
}

// add: 9 a b c
//  assign into <a> the sum of <b> and <c> (podulo 32768)
func add(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, (b+c)%modulo)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, (b+c)%modulo)
	p.index = p.index + 1
}

// and: 12 a b c
//   stores into <a> the bitwise and of <b> and <c>
func and(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b&c)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, b&c)
	p.index = p.index + 1
}

// call: 17 a
//   write the address of the next instruction to the stack and jump to <a>
func call(p *program, r *registers, s *stack) {
	a := p.getNext(r)
	s.push(uint16(p.index) + 1)
	fmt.Fprintf(os.Stderr, "op args: %d, Stack Push: %d", a, p.index+1)
	p.index = int(a)
}

// eq: 4 a b c
//   set <a> to 1 if <b> is equal to <c>;.set it to 0 otherwise
func eq(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)

	set := 0
	if b == c {
		set = 1
	}
	r.set(a, uint16(set))

	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, set)
	p.index = p.index + 1
}

// gt: 5 a b c
//   set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func gt(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)

	set := 0
	if b > c {
		set = 1
	}
	r.set(a, uint16(set))

	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, set)
	p.index = p.index + 1
}

// halt: 0
//   stop execution and terminate the program
func halt(p *program, r *registers, s *stack) {
	fmt.Fprintf(os.Stderr, "op args: n/a")
	os.Exit(0)
}

// in: 20 a
//   read a character from the terminal and write its ascii code to <a>;
//   it can be assumed that once input starts, it will continue until a newline
//   is encountered; this means that you can safely read whole lines from the
//   keyboard and trust that they will be fully read

func in(p *program, r *registers, s *stack) {
	if len(p.input) == 0 {
		if err := p.getChars(); err != nil {
			fmt.Fprintln(os.Stderr, "Error from p.getChars():", err)
			halt(p, r, s)
		}
	}

	a := p.getNextRaw()

	b := p.input[0]
	p.input = p.input[1:]
	r.set(a, b)

	fmt.Fprintf(os.Stderr, "op args: %d, Setting: %d (Char: %s)", a, b, string(b))
	p.index = p.index + 1
}

// jmp: 6 a
//   jump to <a>
func jump(p *program, r *registers, s *stack) {
	p.index = int(p.getNext(r))
	fmt.Fprintf(os.Stderr, "op args: %d, jump location: %d", p.index, p.index)
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
	fmt.Fprintf(os.Stderr, "op args: %d, %d, jump location: %d", a, b, p.index)
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
	fmt.Fprintf(os.Stderr, "op args: %d, %d, jump location: %d", a, b, p.index)
}

// mod: 11 a b c
//   store into <a> the remainder of <b> divided by <c>
func mod(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b%c)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, b%c)
	p.index = p.index + 1
}

// mult: 10 a b c
//   store into <a> the product of <b> and <c> (modulo 32768)
func mult(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, (b*c)%modulo)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, (b*c)%modulo)
	p.index = p.index + 1
}

// noop: 21
//   no operation
func noop(p *program, r *registers, s *stack) {
	p.index = p.index + 1
	fmt.Fprintf(os.Stderr, "op args: n/a")
}

// not: 14 a b
//   stores 15-bit bitwise inverse of <b> in <a>
func not(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	r.set(a, ^b%modulo)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, Setting: %d", a, b, ^b%modulo)
	p.index = p.index + 1
}

// or: 13 a b c
//   stores into <a> the bitwise or of <b> and <c>
func or(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	c := p.getNext(r)
	r.set(a, b|c)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, %d, Setting: %d", a, b, c, b|c)
	p.index = p.index + 1
}

// out: 19 a
//   write the character represented by ascii code <a> to the terminal
func out(p *program, r *registers, s *stack) {
	a := string(p.getNext(r))
	fmt.Print(a)
	fmt.Fprintf(os.Stderr, "op args: %s", a)
	p.index = p.index + 1
}

// push: 2 a
//   push <a> onto the stack
func push(p *program, r *registers, s *stack) {
	a := p.getNext(r)
	s.push(a)
	fmt.Fprintf(os.Stderr, "op args: %d", a)
	p.index = p.index + 1
}

// pop: 3 a
//   remove the top element from the stack and write it into <a>; empty stack = error
func pop(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := s.pop()
	r.set(a, b)
	fmt.Fprintf(os.Stderr, "op args: %d, stack arg: %d", a, b)
	p.index = p.index + 1
}

// ret: 18
//   remove the top element from the stack and jump to it; empty stack = halt
func ret(p *program, r *registers, s *stack) {
	if s.isEmpty() {
		halt(p, r, s)
	}

	a := s.pop()
	fmt.Fprintf(os.Stderr, "op args: n/a, stack arg: %d", a)
	p.index = int(a)
}

// rmem: 15 a b
//   read memory at address <b> and write it to <a>
func rmem(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)
	m := p.memory[b]

	r.set(a, m)
	fmt.Fprintf(os.Stderr, "op args: %d, %d, memory value: %d", a, b, m)
	p.index = p.index + 1
}

// set: 1 a b
//   set register <a> to the value of <b>
func set(p *program, r *registers, s *stack) {
	a := p.getNextRaw()
	b := p.getNext(r)

	if isRegister(a) {
		r.set(a, b)
	}
	fmt.Fprintf(os.Stderr, "op args: %d, %d", a, b)
	p.index = p.index + 1
}

// wmem: 16 a b
//   write the value from <b> into memory at address <a>
func wmem(p *program, r *registers, s *stack) {
	a := p.getNext(r)
	b := p.getNext(r)
	p.memory[a] = b

	fmt.Fprintf(os.Stderr, "op args: %d, %d", a, b)
	p.index = p.index + 1
}
