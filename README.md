# Cago 

Cache Go, key value store in memory

## Table Of Contents

[**Type**](#type)
- [Config](#config)

[**Function**](#function)
- [New(config ...cago.Config)](#how-to-use)
- [Get(key string)](#get)
- [Exist(key string)](#exist)
- [Set(key string, value Compare, maxAge uint64)](#set)
- [Put(key string, value Compare, maxAge uint64)](#put)
- [Remove(key string)](#remove)
- [Clear()](#clear)

## Installation

```sh
go get github.com/jasakode/cago
```

## How to use

```go
package main

import "github.com/jasakode/cago"

func main() {
    err := cago.New(cago.Config{
        Path: "database.db",
    })
    if err != nil {
        panic(err.Error())
    }
}

```
### Config

```go
type Config struct {
	Path                string
	MAX_MEM             uint
	MIN_MEM_ALLOCATION  uint64
	EvictOldestOnMaxMem bool
	TimeoutCheck        uint64
}
```

- **Path** string alamat file untuk menyimpan data
- **MAX_MEM** ukuran maksimal memory yang akan di gunakan
- **MIN_MEM_ALLOCATION** Minimal memory yang akan di bebaskan saat program pertama kali di jalankan
- **EvictOldestOnMaxMem** jika : true, maka cache yang paling pertama di tambahkan akan di hapus jika terjadi kelebihan memory
- **TimeoutCheck** waktu pengecekan data, default 10 detik, maka setiap 10 detik akan melalkukan pengecekan dan penghapusan cahce yang telah kadaluarsa

## Function

### Set
```go
package main

import "github.com/jasakode/cago"

func main() {
    err := cago.Set("name", "Jhon Doe", 10000) // set cache with key name and value jhone in 10 second
    if err != nil {
        panic(err) // if value exist error
    }
    cago.Set("age", 24, 10000) // set cache with key age and value uint 24 in 10 second
    if err != nil {
        panic(err) // if value exist error
    }
    type Person struct {
        Name string `json:"name"`
        Age uint `json:"age"`
    }
    cago.Set("person", Person{ Name: "Jhon Doe", Age: 24 }, 10000) // set cache with key person and value struct Person 24 in 10 second
    if err != nil {
        panic(err) // if value exist error
    }
}

```
### Get

```go
package main

import "github.com/jasakode/cago"

func main() {
    cago.Get[string]("name")
    cago.Get[uint]("age")
    type Person struct {
        Name string `json:"name"`
        Age uint `json:"age"`
    }
    cago.Get[Person]("person")
}
```

### Exist

```go
package main

import "github.com/jasakode/cago"

func main() {
    cago.Exist[string]("name")
}

```

### Put

```go
package main

import "github.com/jasakode/cago"

func main() {
    cago.Put[string]("name", "Cago", 10000) // set or update
}

```

### Remove

```go
package main

import "github.com/jasakode/cago"

func main() {
    cago.Remove[string]("name")
}

```

### Clear

```go
package main

import "github.com/jasakode/cago"

func main() {
    cago.Clear() // clear all cahce
}

```