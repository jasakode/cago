# Cago

[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/2q.svg)](https://pkg.go.dev/github.com/floatdrop/2q)
[![CI](https://github.com/floatdrop/2q/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/2q/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-88.9%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/2q)](https://goreportcard.com/report/github.com/floatdrop/2q)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Thread safe GoLang [Cago](http://www.vldb.org/conf/1994/P439.PDF) cache.

Cache Go, key value store in memory

## Features


```sh
go get github.com/jasakode/cago
```

```go
package main

import "github.com/jasakode/cago"

func main() {
    err := cago.New()
    if err != nil {
        panic(err.Error())
    }
}

```

Execut set in another file in project
```go
package other

import "github.com/jasakode/cago"

func main() {
    store := cago.Get("test")
    fmt.Println(store.Text())
}
```

