// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

/*
	Package cago provides documentation for Go's predeclared identifiers.
	The items documented here are not actually in package cago
	but their descriptions here allow godoc to present documentation
	for the language's special identifiers.
*/

package cago

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jasakode/cago/lib"
	"github.com/jasakode/cago/store"
)

type Config struct {
	// Path ke file database.
	// Ini adalah lokasi spesifik file.
	// Jika path tidak ditentukan, data akan hilang ketika proses dihentikan.
	Path string
	// Memori maksimal yang akan digunakan,
	// ditentukan dalam bit.
	// 8.388.608 bit = 1 MB
	// default : 8589934592 bit / 1 GB
	MAX_MEM uint
}

type App struct {
	mu sync.Mutex
	// Initialize start
	start uint64
	// cago config
	config Config
	// database
	db *database
	// data cache
	data map[string]store.Store
}

var app App = App{}

func New(config ...Config) error {
	if app.config.MAX_MEM == 0 {
		app.config.MAX_MEM = 8388608 * 1204
	}
	app.data = make(map[string]store.Store)
	app.start = uint64(time.Now().UnixMilli())
	app.config = Config{}
	if len(config) > 0 {
		app.config = config[0]
	}
	if app.config.Path != "" {
		if err := app.InitializeDB(); err != nil {
			return err
		}
		if err := app.db.CreateTableIfNotExist(); err != nil {
			return err
		}
		rows, err := app.db.FindALL()
		if err != nil {
			return err
		}
		for i := range *rows {
			val := (*rows)[i]
			// Menambahkan data berdasarkan key tertentu
			app.data[val.Key] = store.ParseStore(val.Value)
		}
		return nil
	}
	return nil
}

// set store
func Set(key string, value store.Compare, maxAge ...uint64) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	if ok {
		return fmt.Errorf("data already exists")
	}
	switch v := any(value).(type) {
	case string:
		data := store.NewStore([]byte(v), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int8:
		data := store.NewStore(lib.Int8ToByte(int8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int16:
		data := store.NewStore(lib.Int16ToByte(int16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int32:
		data := store.NewStore(lib.Int32ToByte(int32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int64:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint8:
		data := store.NewStore(lib.Uint8ToByte(uint8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint16:
		data := store.NewStore(lib.Uint16ToByte(uint16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint32:
		data := store.NewStore(lib.Uint32ToByte(uint32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint64:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case float32, float64:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case any:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Get retrieves a value from the application data.
func Get[K store.Compare](key string) *K {
	app.mu.Lock()
	defer app.mu.Unlock()

	value, ok := app.data[key]
	if !ok {
		return nil // Mengembalikan nil jika key tidak ada
	}

	var result K

	// Menangani setiap tipe dalam switch
	switch any(result).(type) {
	case string:
		result = any(value.Text()).(K)
	case int:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(intValue).(K)
	case int8:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int8:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int8(intValue)).(K) // Konversi jika perlu
	case int16:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int16:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int16(intValue)).(K) // Konversi jika perlu
	case int32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int32(intValue)).(K) // Konversi jika perlu
	case int64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int64(intValue)).(K) // Konversi jika perlu
	case uint:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint(intValue)).(K) // Konversi jika perlu
	case uint8:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint8:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint8(intValue)).(K) // Konversi jika perlu
	case uint16:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint16:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint16(intValue)).(K) // Konversi jika perlu
	case uint32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint32(intValue)).(K) // Konversi jika perlu
	case uint64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint64(intValue)).(K) // Konversi jika perlu
	case float32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving float32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(float32(intValue)).(K) // Konversi jika perlu
	case float64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving float64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(float64(intValue)).(K) // Konversi jika perlu
	default:
		err := value.JSON(&result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return nil // Tangani kesalahan dengan baik
		}
	}

	return &result
}

// Exist checks if a given key exists in the app's data map and returns true if it exists.
func Exist(key string) bool {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	return ok
}

func Put(key string, value store.Compare, maxAge ...uint64) bool {
	app.mu.Lock()
	defer app.mu.Unlock()
	switch v := any(value).(type) {
	case string:
		data := store.NewStore([]byte(v), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int8:
		data := store.NewStore(lib.Int8ToByte(int8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int16:
		data := store.NewStore(lib.Int16ToByte(int16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int32:
		data := store.NewStore(lib.Int32ToByte(int32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case int64:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint8:
		data := store.NewStore(lib.Uint8ToByte(uint8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint16:
		data := store.NewStore(lib.Uint16ToByte(uint16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint32:
		data := store.NewStore(lib.Uint32ToByte(uint32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case uint64:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case float32, float64:
		by, err := json.Marshal(value)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	case any:
		by, err := json.Marshal(value)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
	default:
		fmt.Printf("unsupported type: %T", value)
		return false
	}
	return false
}

func Remove(key string) bool {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	delete(app.data, key)
	if app.db != nil {
		app.db.RemoveByKey(key)
	}
	return ok
}

func Clear() error {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.data = make(map[string]store.Store)
	if app.db != nil {
		return app.db.RemoveAll()
	}
	return nil
}
