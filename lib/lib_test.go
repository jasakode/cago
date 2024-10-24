// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package lib_test

import (
	"testing"

	"github.com/jasakode/cago/lib"
)

// equal membandingkan dua slice byte untuk kesetaraan.
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestUint8ToByte menguji fungsi Uint8ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari uint8 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests digunakan untuk menyimpan input dan output yang diharapkan.
	2. Kasus Uji: Terdapat beberapa kasus uji yang mencakup batas bawah, nilai normal, nilai maksimum untuk tipe uint8, dan nilai boundary.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestUint8ToByte(t *testing.T) {
	tests := []struct {
		input  uint8
		output []byte
	}{
		{0, []byte{0}},     // Kasus batas bawah
		{1, []byte{1}},     // Nilai normal
		{127, []byte{127}}, // Nilai maksimum untuk uint8
		{255, []byte{255}}, // Nilai maksimum untuk uint8 (boundary)
	}
	for _, test := range tests {
		result := lib.Uint8ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Uint8ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestUint16ToByte menguji fungsi Uint16ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari uint16 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan.
	2. Kasus Uji: Mencakup batas bawah, nilai normal, nilai maksimum untuk uint16, dan nilai tengah.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestUint16ToByte(t *testing.T) {
	tests := []struct {
		input  uint16
		output []byte
	}{
		{0, []byte{0, 0}},         // Kasus batas bawah
		{1, []byte{0, 1}},         // Nilai normal
		{255, []byte{0, 255}},     // Nilai maksimum untuk byte pertama
		{256, []byte{1, 0}},       // Nilai maksimum untuk byte kedua
		{65535, []byte{255, 255}}, // Nilai maksimum untuk uint16
		{32768, []byte{128, 0}},   // Nilai tengah
	}

	for _, test := range tests {
		result := lib.Uint16ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Uint16ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestUint32ToByte menguji fungsi Uint32ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari uint32 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup batas bawah, nilai normal, nilai maksimum untuk uint32, dan nilai tengah.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestUint32ToByte(t *testing.T) {
	tests := []struct {
		input  uint32
		output []byte
	}{
		{0, []byte{0, 0, 0, 0}},                  // Kasus batas bawah
		{1, []byte{0, 0, 0, 1}},                  // Nilai normal
		{255, []byte{0, 0, 0, 255}},              // Nilai maksimum untuk byte ketiga
		{256, []byte{0, 0, 1, 0}},                // Nilai dengan byte kedua
		{65535, []byte{0, 0, 255, 255}},          // Nilai maksimum untuk byte kedua dan ketiga
		{4294967295, []byte{255, 255, 255, 255}}, // Nilai maksimum untuk uint32
		{2147483648, []byte{128, 0, 0, 0}},       // Nilai tengah
	}

	for _, test := range tests {
		result := lib.Uint32ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Uint32ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestUint64ToByte menguji fungsi Uint64ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari uint64 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup batas bawah, nilai normal, nilai maksimum untuk uint64, dan nilai tengah.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestUint64ToByte(t *testing.T) {
	tests := []struct {
		input  uint64
		output []byte
	}{
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},                                    // Kasus batas bawah
		{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},                                    // Nilai normal
		{255, []byte{0, 0, 0, 0, 0, 0, 0, 255}},                                // Nilai maksimum untuk byte ketujuh
		{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},                                  // Nilai dengan byte keenam
		{65535, []byte{0, 0, 0, 0, 0, 255, 255, 0}},                            // Nilai maksimum untuk byte keenam dan ketujuh
		{4294967295, []byte{0, 0, 0, 0, 255, 255, 255, 255}},                   // Nilai maksimum untuk uint32
		{18446744073709551615, []byte{255, 255, 255, 255, 255, 255, 255, 255}}, // Nilai maksimum untuk uint64
		{9223372036854775808, []byte{128, 0, 0, 0, 0, 0, 0, 0}},                // Nilai tengah
	}

	for _, test := range tests {
		result := lib.Uint64ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Uint64ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestInt8ToByte menguji fungsi Int8ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari int8 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup nilai batas negatif, nol, nilai positif maksimum, dan nilai lainnya.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestInt8ToByte(t *testing.T) {
	tests := []struct {
		input  int8
		output []byte
	}{
		{-128, []byte{255}}, // Nilai minimum untuk int8
		{-1, []byte{255}},   // Nilai negatif
		{0, []byte{0}},      // Nilai nol
		{1, []byte{1}},      // Nilai positif kecil
		{127, []byte{127}},  // Nilai maksimum untuk int8
	}

	for _, test := range tests {
		result := lib.Int8ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Int8ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestInt16ToByte menguji fungsi Int16ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari int16 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup nilai batas negatif, nol, nilai positif maksimum, dan nilai lainnya.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestInt16ToByte(t *testing.T) {
	tests := []struct {
		input  int16
		output []byte
	}{
		{-32768, []byte{128, 0}},  // Nilai minimum untuk int16
		{-1, []byte{255, 255}},    // Nilai negatif
		{0, []byte{0, 0}},         // Nilai nol
		{1, []byte{0, 1}},         // Nilai positif kecil
		{32767, []byte{127, 255}}, // Nilai maksimum untuk int16
	}

	for _, test := range tests {
		result := lib.Int16ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Int16ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestInt32ToByte menguji fungsi Int32ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari int32 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup nilai batas negatif, nol, nilai positif maksimum, dan nilai lainnya.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestInt32ToByte(t *testing.T) {
	tests := []struct {
		input  int32
		output []byte
	}{
		{-2147483648, []byte{128, 0, 0, 0}},      // Nilai minimum untuk int32
		{-1, []byte{255, 255, 255, 255}},         // Nilai negatif
		{0, []byte{0, 0, 0, 0}},                  // Nilai nol
		{1, []byte{0, 0, 0, 1}},                  // Nilai positif kecil
		{2147483647, []byte{127, 255, 255, 255}}, // Nilai maksimum untuk int32
	}

	for _, test := range tests {
		result := lib.Int32ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Int32ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestInt64ToByte menguji fungsi Int64ToByte dengan berbagai nilai.
// Fungsi ini memeriksa apakah hasil konversi dari int64 ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input dan output yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup nilai batas negatif, nol, nilai positif maksimum, dan nilai lainnya.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestInt64ToByte(t *testing.T) {
	tests := []struct {
		input  int64
		output []byte
	}{
		{-9223372036854775808, []byte{128, 0, 0, 0, 0, 0, 0, 0}},              // Nilai minimum untuk int64
		{-1, []byte{255, 255, 255, 255, 255, 255, 255, 255}},                  // Nilai negatif
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},                                   // Nilai nol
		{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},                                   // Nilai positif kecil
		{9223372036854775807, []byte{127, 255, 255, 255, 255, 255, 255, 255}}, // Nilai maksimum untuk int64
	}

	for _, test := range tests {
		result := lib.Int64ToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("Int64ToByte(%d) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestStringToByte menguji fungsi StringToByte dengan berbagai nilai string.
// Fungsi ini memeriksa apakah hasil konversi dari string ke []byte sesuai dengan yang diharapkan.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input string dan output []byte yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup string kosong, string normal, dan string dengan karakter spesial.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestStringToByte(t *testing.T) {
	tests := []struct {
		input  string
		output []byte
	}{
		{"", []byte{}}, // Kasus string kosong
		{"hello", []byte{'h', 'e', 'l', 'l', 'o'}},                               // String normal
		{"12345", []byte{'1', '2', '3', '4', '5'}},                               // String dengan angka
		{"!@#$%^&*()", []byte{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')'}}, // String dengan karakter spesial
	}

	for _, test := range tests {
		result := lib.StringToByte(test.input)
		if !equal(result, test.output) {
			t.Errorf("StringToByte(%q) = %v; expected %v", test.input, result, test.output)
		}
	}
}

// TestStringToByteASCII menguji fungsi StringToByteASCII dengan berbagai nilai string.
// Fungsi ini memeriksa apakah hasil konversi dari string ke []byte sesuai dengan yang diharapkan,
// terutama dalam menangani karakter ASCII dan non-ASCII.
/*
	1. Test Structure: Struktur tests berisi kombinasi nilai input string dan output []byte yang diharapkan untuk pengujian.
	2. Kasus Uji: Mencakup string kosong, string dengan karakter ASCII, dan string dengan karakter non-ASCII.
	3. Comparing Results: Fungsi equal digunakan untuk membandingkan dua slice byte, memastikan hasilnya sesuai dengan yang diharapkan.
*/
func TestStringToByteASCII(t *testing.T) {
	tests := []struct {
		input  string
		output []byte
	}{
		{"", []byte{}}, // Kasus string kosong
		{"hello", []byte{'h', 'e', 'l', 'l', 'o'}},                                                             // String dengan karakter ASCII
		{"12345", []byte{'1', '2', '3', '4', '5'}},                                                             // String dengan angka
		{"hello, 世界", []byte{'h', 'e', 'l', 'l', 'o', ',', 0}},                                                 // Kombinasi ASCII dan non-ASCII
		{"ASCII: !@#$%^&*", []byte{'A', 'S', 'C', 'I', 'I', ':', ' ', '!', '@', '#', '$', '%', '^', '&', '*'}}, // String dengan karakter spesial
		{"Café", []byte{'C', 'a', 'f', 'e', 0}},                                                                // Contoh dengan karakter non-ASCII 'é'
	}

	for _, test := range tests {
		result := lib.StringToByteASCII(test.input)
		if !equal(result, test.output) {
			t.Errorf("StringToByteASCII(%q) = %v; expected %v", test.input, result, test.output)
		}
	}
}
