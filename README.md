# Cago

[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/2q.svg)](https://pkg.go.dev/github.com/floatdrop/2q)
[![CI](https://github.com/floatdrop/2q/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/2q/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-88.9%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/2q)](https://goreportcard.com/report/github.com/floatdrop/2q)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Thread safe GoLang [Cago](http://www.vldb.org/conf/1994/P439.PDF) cache.

Cache Go, key value store in memory



### Value structure representation

- 5 bytes signature / 63 67 65 71 79
- 8 bytes start date
- chunk | 
    Length (4 bytes) | 
    STX start of text (1 byte) |
    ETX end of text (1 byte) |
    Data (variable) | 
    Create At 8 bytes | 
    Update At 8 bytes | 
    GS group sparator (1 byte) |



