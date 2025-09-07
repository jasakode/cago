// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package store_test

import (
	"testing"
	"time"

	"github.com/jasakode/cago/src/lib"
	"github.com/jasakode/cago/src/store"
)

const (
	CreateAtIndex  = 0  // Indeks untuk waktu pembuatan dalam penyimpanan
	UpdateAtIndex  = 8  // Indeks untuk waktu pembaruan dalam penyimpanan
	MaxAgeIndex    = 16 // Indeks untuk usia maksimum data dalam penyimpanan
	LengthIndex    = 24 // Indeks untuk panjang data yang disimpan
	DataStartIndex = 32 // Indeks awal untuk data aktual dalam penyimpanan
)

// TestNewStore menguji fungsi NewStore dengan berbagai nilai data dan maxAge.
// Fungsi ini memastikan bahwa Store yang dihasilkan memiliki metadata dan data yang sesuai.
/*
	1. Test Structure: Struktur tests digunakan untuk menyimpan input dan output yang diharapkan.
	2. Kasus Uji: Terdapat kasus uji dengan data normal dan usia maksimum untuk memverifikasi pembuatan Store.
	3. Comparing Results: Memastikan panjang Store, timestamp, nilai max age, dan data yang disalin sesuai dengan yang diharapkan.
*/
func TestNewStore(t *testing.T) {
	data := []byte("example data")
	maxAge := uint64(60) // Usia maksimum dalam detik

	// Buat Store baru

	s := store.NewStore(data, maxAge)

	// Pastikan panjang Store sesuai
	expectedLength := DataStartIndex + len(data)
	if len(s) != expectedLength {
		t.Errorf("expected length %d, got %d", expectedLength, len(s))
	}

	// Verifikasi timestamp creation
	expectedCreateAt := uint64(time.Now().UnixMilli())
	createAt := s.CreateAt()
	if createAt < expectedCreateAt-1000 || createAt > expectedCreateAt {
		t.Errorf("CreateAt out of range: expected ~%d, got %d", expectedCreateAt, createAt)
	}

	// Verifikasi nilai max age
	if s.MaxAge() != maxAge {
		t.Errorf("expected max age %d, got %d", maxAge, s.MaxAge())
	}

	// Verifikasi panjang data
	if s.Length() != uint64(len(data)) {
		t.Errorf("expected length %d, got %d", len(data), s.Length())
	}

	// Verifikasi data yang disalin
	if string(s.Bytes()) != string(data) {
		t.Errorf("expected data %s, got %s", data, s.Bytes())
	}
}

// TestParseStore menguji fungsi ParseStore untuk memverifikasi penguraian data menjadi Store yang valid.
// Fungsi ini memastikan bahwa data valid dapat diparsing dengan benar, dan data tidak valid mengembalikan Store kosong.
/*
	1. Test Structure: Struktur tests digunakan untuk menguji kasus dengan data valid dan tidak valid.
	2. Kasus Uji: Menggunakan data yang cukup panjang untuk memuat semua metadata dan data yang disimpan.
	3. Validasi Output: Memastikan Store yang dihasilkan dari data valid tidak kosong dan Store yang dihasilkan dari data tidak valid adalah kosong.
*/
func TestParseStore(t *testing.T) {
	// Kasus uji dengan data yang valid
	validData := make([]byte, DataStartIndex+8) // Cukup panjang untuk metadata
	copy(validData[CreateAtIndex:UpdateAtIndex], lib.Uint64ToByte(uint64(time.Now().UnixMilli())))
	copy(validData[UpdateAtIndex:MaxAgeIndex], make([]byte, 8))
	copy(validData[MaxAgeIndex:LengthIndex], lib.Uint64ToByte(60))
	copy(validData[LengthIndex:], lib.Uint64ToByte(8))
	copy(validData[DataStartIndex:], []byte("data"))

	// Mengurai Store dari data valid
	s := store.ParseStore(validData)

	// Pastikan Store tidak kosong
	if len(s) == 0 {
		t.Error("expected non-empty Store, got empty")
	}

	// Kasus uji dengan data tidak valid
	invalidData := []byte("invalid")
	storeInvalid := store.ParseStore(invalidData)

	// Pastikan Store kosong untuk data tidak valid
	if len(storeInvalid) != 0 {
		t.Error("expected empty Store for invalid data, got non-empty")
	}
}
