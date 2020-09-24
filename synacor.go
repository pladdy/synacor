package synacor

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// Machine which represents a program (in memory) that can be run.
type Machine struct {
	Program   *program
	Stack     *stack
	Registers *registers
}

// NewMachine returns a new Machine type.
func NewMachine() Machine {
	return Machine{&program{}, &stack{}, &registers{}}
}

// HasMoreOps returns true if Machine has more operations to run.
func (m Machine) HasMoreOps() bool {
	if m.Program.index < len(m.Program.memory) {
		return true
	}
	return false
}

// Load takes a path to a binary and loads it into the Machine
func (m Machine) Load(s string) {
	m.Program.load(s)
}

// NextOp returns the
//   - name of the next operation
//   - code of the next operation
//   - arguments for the next operation
func (m Machine) NextOp() (name string, opCode uint16, args []uint16) {
	p := m.Program
	oc := opcode(p.memory[p.index])
	properties := operatorPropertyMap[oc]

	for i := 0; i < properties.args; i++ {
		args = append(args, p.getNextRaw())
	}
	p.index = p.index + 1

	return properties.name, uint16(oc), args
}

// Run the loaded program.
func (m Machine) Run() {
	// hacks
	const hackSetReg = 6034
	hackedSetReg := false

	const hackEndOfRegisterTests = 521
	hackedEndOfRegisterTests := false

	const hackCallAcker = 6027
	hackedCallAcker := false

	p := m.Program
	for p.index < len(p.memory) {

		if p.index > hackEndOfRegisterTests && hackedEndOfRegisterTests == false {
			// run `make teleporter` and result is entered here
			m.Registers.set(registerEnd, uint16(25734))
			hackedEndOfRegisterTests = true
		}

		// Below hacks will allow you to get to next part of the program

		// Disable check to see program calls to see what value is being checked.
		// The check takes too long to allow to run, so it has to be disabled.
		if p.index >= hackCallAcker && hackedCallAcker == false {
			m.Registers.set(registerStart, 0)
			m.Registers.set(registerStart+1, 0)
			hackedCallAcker = true
		}

		// Set first register to expected final result (6).  If the value in
		// register 8 is wrong, the code you get will be invalid.
		if p.index == hackSetReg && hackedSetReg == false {
			m.Registers.set(registerStart, 6)
			hackedSetReg = true
		}

		v := opcode(p.memory[p.index])
		ops := operatorPropertyMap[v]

		fmt.Fprintf(os.Stderr, "%d %s (%d) ", p.index, ops.name, v)

		operatorFunctionMap[v](p, m.Registers, m.Stack)

		fmt.Fprintf(os.Stderr, " Stack: %d, Registers: %d", m.Stack, m.Registers)
		fmt.Fprintf(os.Stderr, " Input: '%s'\n", inputToString(p.input))

		// custom debug statements
		// first char typed into stdin gets set in a register...
		if strings.Contains(inputToString(p.input), "se teleporter") {
			fmt.Fprintln(os.Stderr, "  'use teleporter' called")
		}
	}
}

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

// inputToString can only return what's in the input property; the first char
// of any input is saved to a register and the rest of 'input' is read as in
// operations are called.
func inputToString(input []uint16) string {
	var b strings.Builder

	for i := 0; i < len(input); i++ {
		fmt.Fprintf(&b, "%s", string(rune((input[i]))))
	}
	return b.String()
}

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
