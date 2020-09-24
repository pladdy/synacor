package main

import "testing"

var result int

func BenchmarkAck(b *testing.B) {
	var r int

	for n := 0; n < b.N; n++ {
		r = ack(2, 10000)
	}
	result = r
}

func BenchmarkAckRecursive(b *testing.B) {
	var r int

	for n := 0; n < b.N; n++ {
		r = ackRecursive(2, 10000)
	}

	result = r
}
