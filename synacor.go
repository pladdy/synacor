package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

var registers [8]uint16
var stack []uint16

type memoryStack []uint16

func (m *memoryStack) pop() uint16 {
	last := len(*m) - 1
	value := (*m)[last]
	*m = (*m)[:last]
	return value
}

func (m *memoryStack) push(i uint16) {
	*m = append(*m, i)
	return
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

func loadProgram(file string) []uint16 {
	fh, err := os.Open("./challenge.bin")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(fh)
	memory := memoryStack{}
	i := 0
	for {
		le, err := readNext(reader)
		if err == io.EOF {
			break
		}

		// memory = append(memory, le)
		memory.push(le)
		i = i + 1
	}

	if err := fh.Close(); err != nil {
		panic(err)
	}

	return memory
}

func readNext(reader io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func main() {
	memory := loadProgram("./challenge.bin")
	fmt.Println("Program loaded into memory.")

	for i, value := range memory {
		// i = handleOpcode(value, i, &memory)
		fmt.Printf("Memory index: %d, Decimal: %d, Binary: %b\n", i, value, value)
	}
}
