package cago_test

import (
	"bytes"
	"testing"
)

func BenchmarkCompareString(b *testing.B) {
	str1 := "hello world"
	str2 := "hello world"
	for i := 0; i < b.N; i++ {
		if str1 == str2 {
			continue
		}
	}
}

func BenchmarkCompareByte(b *testing.B) {
	byte1 := []byte("hello world")
	byte2 := []byte("hello world")
	for i := 0; i < b.N; i++ {
		if bytes.Equal(byte1, byte2) {
			continue
		}
	}
}
