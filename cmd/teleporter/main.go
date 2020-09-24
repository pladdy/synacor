package main

import "fmt"

const vmMaxInt = 32768

// someFunc because I don't know why we use it...
func someFunc(c int) int {
	return (2*c + 1 + (c+1)*c) % vmMaxInt
}

func main() {
	// fmt.Println(someFunc(0))
	// fmt.Println(someFunc(1))
	// fmt.Println(someFunc(2))

	for c := 1; c < vmMaxInt; c++ {
		var values [vmMaxInt]int
		values[0] = someFunc(c)

		for i := 1; i < vmMaxInt; i++ {
			values[i] = (values[i-1]*(c+1) + 1 + 2*c) % vmMaxInt
		}

		value := values[values[c]]
		fmt.Println("C:", c, "Value[c]:", values[c], "Value:", value)
		if value == 6 {
			fmt.Println("  Answer:", value)
			break
		}
	}
}

/*
  Below is what I started with but couldn't calculate the values fast enough.
	Instead I'm trying to use https://github.com/pankdm/synacor-challenge/blob/master/teleport.cpp
	which seems to do this way faster via math instead of the garbage I wrote.
*/
var memoize = make(map[int]map[int]int)

// The recurisve function runs faster than the iterative, but the program as a
// whole is too slow.
func ackRecursive(m, n int) int {
	//fmt.Println("Recursive, m:", m, "n:", n, "size:")

	// check memoize first
	if _, ok := memoize[m][n]; ok {
		return memoize[m][n]
	}

	if m == 0 {
		return (n + 1)
	}

	if m > 0 && n == 0 {
		return ackRecursive(m-1, 1)
	}

	return ackRecursive(m-1, ackRecursive(m, n-1))
}

// The iterative version isn't faster than recursion, and the program is too
// slow still.
func ack(m, n int) int {
	value := map[int]int{n: 0}
	size := 0

	for {
		// fmt.Println("Iterative m:", m, "n:", n, "size:", size, "value:", value)

		if m == 0 {
			n++

			if size == 0 {
				break
			}

			size--
			m = value[size]
			continue
		}

		if n == 0 {
			m--
			n = 1
			continue
		}

		value[size] = m - 1
		size++
		n--
	}

	return n
}

func main_too_slow() {
	fmt.Println("Running ackermann function...")

	// fmt.Println(ackRecursive(3, 12))
	// fmt.Println(ack(3, 12))

	for m := 0; m <= 3; m++ {
		memoize[m] = make(map[int]int)
		//
		for n := 0; n <= vmMaxInt; n++ {
			fmt.Println("Computing m:", m, "n:", n)
			result := ackRecursive(m, n)
			memoize[m][n] = result
			fmt.Println("Raw result for m:", m, "n:", n, ":", result, result%vmMaxInt)
		}
		// fmt.Println("cache:", memoize)
	}
}
