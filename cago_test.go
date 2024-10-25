// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago_test

import (
	"bytes"
	"fmt"
	"sync"
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
	var wg sync.WaitGroup

	// Hapus file database sebelum test dimulai

	t.Cleanup(func() {
		// os.Remove("db.db")
	})

	// Inisialisasi cago dengan database baru
	cago.New(cago.Config{
		Path: "db.db",
	})

	// Tambahkan WaitGroup untuk menunggu proses selesai
	wg.Add(1)

	// Luncurkan goroutine untuk menunggu 3 detik
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
	}()

	// Tunggu sampai goroutine selesai
	wg.Wait()

	// Test untuk key "jhon"
	rs := cago.Get[string]("jhon")
	if rs != nil {
		fmt.Println(*rs)
	} else {
		fmt.Println("Data Not Found !!!")
	}

	// cago.Set("hello", "HALLO KAMU", 5000)
	// cago.Set("jhon", "HALLO KAMU Jhon", 60000 * 60)
	// time.Sleep(1 * time.Second)
	// fmt.Println(cago.Size())
	// fmt.Println(*cago.Get[string]("hello"))
	// fmt.Println(*cago.Get[string]("jhon"))
	// time.Sleep(23 * time.Second)
	// fmt.Println(cago.Size())
	// fmt.Println(cago.Get[string]("hello"), cago.Get[string]("jhon"))
}
