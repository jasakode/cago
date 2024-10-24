// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago_test

import (
	"testing"

	"github.com/jasakode/cago"
)

func TestDbConnection(t *testing.T) {
	t.Cleanup(func() {})
	err := cago.New(cago.Config{
		Path: "db.db",
	})
	if err != nil {
		t.Fail()
	}
}
