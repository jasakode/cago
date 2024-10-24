// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/jasakode/cago"
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

type Person struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func TestApp(t *testing.T) {
	t.Cleanup(func() {})

	cago.New()
	time.Sleep(1 * time.Second)
	cago.Set("hello", uint64(2327632839))
	fmt.Println(*cago.Get[uint64]("hello"))
}
