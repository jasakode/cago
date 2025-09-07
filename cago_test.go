// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago_test

import (
	"testing"
	"time"

	"github.com/jasakode/cago"
)

func setup()   { _ = cago.New(cago.Config{CleanInterval: 20 * time.Millisecond}) }
func destroy() { cago.Close() }

func TestSetGetAndExpire(t *testing.T) {
	setup()
	defer destroy()

	if err := cago.Set("greeting", "hello", 80*time.Millisecond); err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}

	if v, ok := cago.Get[string]("greeting"); !ok || v != "hello" {
		t.Fatalf("Get expected 'hello', got %q ok=%v", v, ok)
	}

	time.Sleep(120 * time.Millisecond)
	if _, ok := cago.Get[string]("greeting"); ok {
		t.Fatalf("expected key to be expired")
	}
	if cago.Exist("greeting") {
		t.Fatalf("Exist should be false after expiration")
	}
}

func TestSetConflictAndPut(t *testing.T) {
	setup()
	defer destroy()

	if err := cago.Set("k", 123, 0); err != nil {
		t.Fatalf("unexpected error on first Set: %v", err)
	}
	if err := cago.Set("k", 456, 0); err == nil {
		t.Fatalf("expected ErrKeyExists on second Set")
	}

	cago.Put("k", 456, 0)
	if v, ok := cago.Get[int]("k"); !ok || v != 456 {
		t.Fatalf("Put did not overwrite value: got %v ok=%v", v, ok)
	}
}

func TestRemoveAndClear(t *testing.T) {
	setup()
	defer destroy()

	cago.Put("a", "x", 0)
	cago.Put("b", "y", 0)

	if ok := cago.Remove("a"); !ok {
		t.Fatalf("expected Remove to return true")
	}
	if _, ok := cago.Get[string]("a"); ok {
		t.Fatalf("expected 'a' to be removed")
	}

	cago.Clear()
	if cago.Exist("b") {
		t.Fatalf("expected 'b' not to exist after Clear")
	}
}
