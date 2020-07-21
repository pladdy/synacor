package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const modulo = 32768 // 2 << 14
const register0 = 32768
const register7 = 32775

var memory [2 << 14]uint16
var registers [8]uint16
var stack []uint16

func isValid(u uint16) bool {
	return u <= register7
}

func isLiteralValue(u uint16) bool {
	return u < modulo
}

func isRegister(u uint16) bool {
	return u >= register0 && u <= register7
}

func main() {
	fmt.Println("Virtual Machine goes here")

	fh, err := os.Open("./challenge.bin")
	defer fh.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(fh)
	i := 0
	for {
		buf := make([]byte, 2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		le := binary.LittleEndian.Uint16(buf)
		fmt.Printf("Decimal: %d, Binary: %b\n", le, le)
		memory[i] = le
	}
}
