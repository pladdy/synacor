package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const maxMemory = 32767         // var memory [2 << 14]uint16
const maxAllowedLiteral = 32767 // 2 << 14 - 1
const modulo = 32768
const registerStart = 32768
const registerEnd = 32775

const (
	register0 = iota + registerStart
	register1
	register2
	register3
	register4
	register5
	register6
	register7
)

type program struct {
	index  int
	memory []uint16
}

// This returns the value and shifts the provided index
func (p *program) getNext(r *registers) uint16 {
	p.index = p.index + 1
	value := p.memory[p.index]

	// what if the index is a register we want to set a value to?
	if isRegister(value) {
		value = r.Get(value)
	}
	return value
}

// Similar to getNext, but returns the raw value, not the value in the register
// if the value is a register.
func (p *program) getNextRaw() uint16 {
	p.index = p.index + 1
	return p.memory[p.index]
}

func (p *program) load(file string) {
	fh, err := os.Open(filepath.Clean(file))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(fh)
	i := 0
	for {
		le, err := readNext(reader)
		if err == io.EOF {
			break
		}

		p.memory = append(p.memory, le)
		i = i + 1
	}

	if err := fh.Close(); err != nil {
		panic(err)
	}
}

type registers [8]uint16

func (r *registers) Get(register uint16) uint16 {
	return r[register%registerStart]
}

func (r *registers) Set(register uint16, value uint16) {
	r[register%registerStart] = value
}

type stack []uint16

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

func (s *stack) Pop() uint16 {
	if s.isEmpty() {
		return 0
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value
}

func (s *stack) Push(v uint16) {
	*s = append(*s, v)
}

/* helpers */

func isValid(u uint16) bool {
	return u <= registerEnd
}

func isLiteralValue(u uint16) bool {
	return u <= maxAllowedLiteral
}

func isRegister(u uint16) bool {
	return u >= register0 && u <= register7
}

func readNext(reader io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func main() {
	r := registers{}
	p := program{}
	p.load("./challenge.bin")
	fmt.Println("Program loaded into memory.")

	for p.index < len(p.memory) {
		v := p.memory[p.index]
		fmt.Printf("DEBUG: Memory index: %d, Decimal: %d, Binary: %b\n", p.index, v, v)
		operatorMap[opcode(v)](&p, &r)
	}
}
