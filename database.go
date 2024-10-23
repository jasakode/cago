// Copyright 2024 The Cago Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/*
	Package builtin provides documentation for Go's predeclared identifiers.
	The items documented here are not actually in package builtin
	but their descriptions here allow godoc to present documentation
	for the language's special identifiers.
*/

package cago

type Database struct {
}

var db *Database

func NewDB() *Database {
	d := Database{}
	db = &d
	return db
}

func (d *Database) FindALL(offset int, limit int) {

}

func (d *Database) Find() {}

func (d *Database) Remove(key string) {}
