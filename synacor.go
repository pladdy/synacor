package synacor

import (
	"bufio"
	"encoding/binary"
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
	input  []uint16
}

// This returns the value and shifts the provided index
func (p *program) getNext(r *registers) uint16 {
	value := p.getNextRaw()

	// what if the index is a register we want to set a value to?
	if isRegister(value) {
		value = r.get(value)
	}
	return value
}

// Similar to getNext, but returns the raw value, not the value in the register
// if the value is a register.
func (p *program) getNextRaw() uint16 {
	p.index = p.index + 1
	// fmt.Println("  next index:", p.index, "value:", p.memory[p.index])
	return p.memory[p.index]
}

func (p *program) getChars() error {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	for _, c := range input {
		p.input = append(p.input, uint16(c))
	}
	// fmt.Println("Input captured")
	return nil
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

func (r *registers) get(register uint16) uint16 {
	return r[register%registerStart]
}

func (r *registers) set(register uint16, value uint16) {
	if isRegister(value) {
		// fmt.Println("Value is a register:", value)
		os.Exit(1)
	}
	r[register%registerStart] = value
}

type stack []uint16

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

func (s *stack) pop() uint16 {
	if s.isEmpty() {
		return 0
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value
}

func (s *stack) push(v uint16) {
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
