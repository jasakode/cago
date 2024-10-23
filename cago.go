// Copyright 2024 The Cago Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package cago provides documentation for Go's predeclared identifiers.
	The items documented here are not actually in package cago
	but their descriptions here allow godoc to present documentation
	for the language's special identifiers.
*/

package cago

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type Config struct {
	// Path ke file database.
	// Ini adalah lokasi spesifik file.
	// Jika path tidak ditentukan, data akan hilang ketika proses dihentikan.
	path string
	// Memori maksimal yang akan digunakan,
	// ditentukan dalam bit.
	// 8.388.608 bit = 1 MB
	// default : 1 GB
	MAX_MEM uint
}

type App struct {
	mu     sync.Mutex
	config Config
	data   []byte
}

func New[T string | Config](config ...T) *App {
	app := App{}
	app.config = Config{}
	if len(config) > 0 {
		switch v := any(config[0]).(type) {
		case string:
			app.config.path = v
		case Config:
			app.config = v
		}
	}
	app.init()
	return &app
}

func (app *App) init() {
	app.mu.Lock()
	defer app.mu.Unlock()
	if app.config.MAX_MEM == 0 {
		app.config.MAX_MEM = 8388608 * 1204
	}
	app.data = make([]byte, 13)
	copy(app.data[0:5], []byte{63, 67, 65, 71, 79})
	timestamp := time.Now().UnixMilli()
	fmt.Println(uint64(timestamp))
	binary.BigEndian.PutUint64(app.data[4:12], uint64(timestamp))
}

func ToByte(inp int32) []byte {
	m := make([]byte, 4)
	binary.BigEndian.PutUint32(m, uint32(inp))
	return m
}

func GetMapSize(m map[string]interface{}) (size int, count int) {
	// Ukuran pointer dari map
	size = int(unsafe.Sizeof(m))
	// Menghitung jumlah elemen
	for range m {
		count++
	}
	return size, count
}

// ln, err := linux.NewSecureKey("")
// fmt.Println()
// var emoji rune = '😊'
// fmt.Println(ToByte(emoji), emoji)
// fmt.Println(ToByte(-1), -1)
// fmt.Println(ToByte(256), 256)

// app := New[string]()
// // fmt.Println(string(app.data))
// fmt.Printf("%b\n", app.data[4:12])
// fmt.Printf("%d\n", binary.BigEndian.Uint64(app.data[4:12]))
