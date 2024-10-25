// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jasakode/cago"
)

func TestDbConnection(t *testing.T) {
	t.Cleanup(func() {})
	err := cago.New()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	time.Sleep(2 * time.Second)
	if err := cago.Set("jhon", "HALLO KAMU Jhon", 60000*60); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	time.Sleep(1 * time.Second)
	// Test untuk key "jhon"
	rs := cago.Get[string]("jhon")
	if rs != nil {
		fmt.Println(*rs)
	} else {
		fmt.Println("Data Not Found !!!")
	}

	cago.Put("jhon", "Babi Kau !!!")

	time.Sleep(1 * time.Second)
	// Test untuk key "jhon"
	rss := cago.Get[string]("jhon")
	if rs != nil {
		fmt.Println(*rss)
	} else {
		fmt.Println("Data Not Found !!!")
	}
}
