// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jasakode/cago/lib"
)

// Store adalah tipe data yang merepresentasikan sekumpulan byte.
// Tipe ini dapat digunakan untuk menyimpan data biner dalam bentuk slice byte.
type Store []byte

// Compare adalah interface yang mendefinisikan tipe data yang dapat dibandingkan.
// Interface ini mencakup berbagai tipe numerik, string, dan tipe data lainnya yang dapat digunakan dalam operasi perbandingan.
// Tipe yang diizinkan mencakup uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64,
// serta tipe dasar int, uint, string, dan any.
type Compare interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 | int | uint | string | any
}

const (
	CreateAtIndex  = 0  // Indeks untuk waktu pembuatan dalam penyimpanan
	UpdateAtIndex  = 8  // Indeks untuk waktu pembaruan dalam penyimpanan
	MaxAgeIndex    = 16 // Indeks untuk usia maksimum data dalam penyimpanan
	LengthIndex    = 24 // Indeks untuk panjang data yang disimpan
	DataStartIndex = 32 // Indeks awal untuk data aktual dalam penyimpanan
)

// NewStore membuat penyimpanan baru dengan metadata dan data yang diberikan.
// Fungsi ini menginisialisasi struktur penyimpanan dengan waktu pembuatan,
// waktu pembaruan (default ke nol), usia maksimum, panjang data, dan data aktual.
//
// Parameter:
// - data: Data biner yang akan disimpan.
// - maxAge: Usia maksimum yang diperbolehkan untuk data (opsional).
//
// Mengembalikan:
// - Store: Struktur penyimpanan yang berisi metadata dan data yang diberikan.
func NewStore(data []byte, maxAge ...uint64) Store {
	MaxAge := uint64(0) // Inisialisasi usia maksimum ke nol
	if len(maxAge) > 0 {
		MaxAge = maxAge[0] // Jika ada argumen maxAge, ambil nilainya
	}

	// Membuat slice Store dengan panjang yang cukup untuk metadata dan data
	s := make(Store, DataStartIndex+len(data))
	copy(s[CreateAtIndex:UpdateAtIndex], lib.Uint64ToByte(uint64(time.Now().UnixMilli()))) // Menyimpan waktu pembuatan
	copy(s[UpdateAtIndex:MaxAgeIndex], make([]byte, 8))                                    // Menyimpan nilai nol untuk waktu pembaruan
	copy(s[MaxAgeIndex:LengthIndex], lib.Uint64ToByte(MaxAge))                             // Menyimpan usia maksimum
	copy(s[LengthIndex:], lib.Uint64ToByte(uint64(len(data))))                             // Menyimpan panjang data
	copy(s[DataStartIndex:], data)                                                         // Menyalin data aktual setelah metadata
	return s                                                                               // Mengembalikan struktur penyimpanan yang telah dibuat
}

// ParseStore menguraikan data byte dan mengembalikan Store yang sesuai.
// Fungsi ini memastikan bahwa data memiliki panjang yang cukup untuk
// mencakup semua metadata yang diperlukan sebelum mengembalikannya.
//
// Parameter:
// - data: Data biner yang akan diuraikan menjadi Store.
//
// Mengembalikan:
// - Store: Struktur penyimpanan yang berisi metadata dan data yang diberikan.
// - Jika data tidak valid, kembalikan Store kosong.
func ParseStore(data []byte) Store {
	// Pastikan panjang data cukup untuk menampung semua metadata
	if len(data) < DataStartIndex {
		return Store{} // Mengembalikan Store kosong jika data tidak valid
	}

	return Store(data) // Mengembalikan data sebagai Store
}

// Values mengembalikan seluruh data yang disimpan dalam Store sebagai slice byte.
// Fungsi ini mengakses nilai yang disimpan di dalam Store dan mengembalikannya
// tanpa memodifikasi data.
//
// Mengembalikan:
//   - []byte: Data yang tersimpan dalam Store dalam bentuk slice byte.
func (s Store) Values() []byte {
	return s
}

// CreateAt mengembalikan timestamp saat store dibuat.
// Fungsi ini mengambil nilai timestamp dari indeks yang ditentukan dalam
// struktur Store. Timestamp ini disimpan dalam format big-endian
// di dalam byte slice `s` pada rentang indeks dari CreateAtIndex
// hingga UpdateAtIndex.
//
// Mengembalikan:
//   - uint64: Timestamp dalam format Unix yang menunjukkan waktu pembuatan
//     dari store dalam milidetik.
func (s Store) CreateAt() uint64 {
	return binary.BigEndian.Uint64(s[CreateAtIndex:UpdateAtIndex])
}

// UpdateAt mengembalikan timestamp terakhir kali store diperbarui.
// Fungsi ini mengambil nilai timestamp dari indeks yang ditentukan dalam
// struktur Store. Timestamp ini disimpan dalam format big-endian
// di dalam byte slice `s` pada rentang indeks dari UpdateAtIndex
// hingga MaxAgeIndex.
//
// Mengembalikan:
//   - uint64: Timestamp dalam format Unix yang menunjukkan waktu terakhir
//     pembaruan dari store dalam milidetik. Nilai ini akan bernilai nol
//     jika store belum pernah diperbarui.
func (s Store) UpdateAt() uint64 {
	return binary.BigEndian.Uint64(s[UpdateAtIndex:MaxAgeIndex])
}

// SetUpdateAt menetapkan timestamp terakhir kali store diperbarui.
// Fungsi ini menerima parameter `date` yang merupakan timestamp dalam
// format Unix dan mengupdate nilai timestamp di dalam store pada
// indeks yang ditentukan (UpdateAtIndex hingga MaxAgeIndex).
//
// Parameter:
//   - date (uint64): Timestamp dalam format Unix yang menunjukkan waktu
//     saat store diperbarui.
//
// Mengembalikan:
//   - Store: Mengembalikan instance Store yang telah diperbarui
//     dengan timestamp baru.
func (s Store) SetUpdateAt(date uint64) Store {
	binary.BigEndian.PutUint64(s[UpdateAtIndex:MaxAgeIndex], date)
	return s
}

// Length mengembalikan panjang data yang disimpan dalam store.
// Jika parameter opsional `all` diisi dan bernilai true, maka
// panjang keseluruhan store akan dikembalikan. Jika tidak,
// fungsi ini akan membaca nilai panjang dari indeks yang ditentukan
// (LengthIndex) dan mengembalikannya sebagai uint64.
//
// Parameter:
// - all (opsional): Jika diisi true, mengembalikan panjang seluruh store.
//
// Mengembalikan:
// - uint64: Panjang data yang disimpan atau panjang keseluruhan store jika all true.
func (s Store) Length(all ...bool) uint64 {
	if len(all) > 0 && all[0] {
		return uint64(len(s))
	}
	return binary.BigEndian.Uint64(s[LengthIndex:])
}

// MaxAge mengembalikan usia maksimum yang disimpan dalam store.
// Fungsi ini mengambil 8 byte dari penyimpanan, dimulai dari
// indeks MaxAgeIndex dan mengonversinya menjadi uint64.
//
// Mengembalikan:
//   - uint64: Usia maksimum yang disimpan dalam store.
func (s Store) MaxAge() uint64 {
	return binary.BigEndian.Uint64(s[MaxAgeIndex:LengthIndex])
}

// SetMaxAge mengatur usia maksimum yang disimpan dalam store.
// Fungsi ini menerima nilai maxAge sebagai parameter dan menyimpannya
// dalam penyimpanan mulai dari indeks MaxAgeIndex. Jika panjang
// data tidak mencukupi untuk menyimpan usia maksimum, fungsi ini
// akan mengembalikan kesalahan.
//
// Parameter:
//   - maxAge: Usia maksimum yang ingin diatur dalam store.
//
// Mengembalikan:
//   - Store: Struktur penyimpanan yang diperbarui dengan usia maksimum baru.
func (s Store) SetMaxAge(maxAge uint64) Store {
	// Mengonversi maxAge ke byte dan menyimpannya di penyimpanan
	copy(s[MaxAgeIndex:LengthIndex], lib.Uint64ToByte(maxAge))
	return s // Mengembalikan struktur penyimpanan yang telah diperbarui
}

// SetLength menetapkan panjang data yang disimpan dalam store.
// Fungsi ini menerima parameter `length` yang merupakan panjang data
// yang ingin disimpan, dan mengupdate nilai panjang di dalam store
// pada indeks yang ditentukan (LengthIndex).
//
// Parameter:
// - length (uint64): Panjang data yang akan disimpan di dalam store.
//
// Mengembalikan:
//   - Store: Mengembalikan instance Store yang telah diperbarui dengan
//     panjang data baru.
func (s Store) SetLength(length uint64) Store {
	binary.BigEndian.PutUint64(s[LengthIndex:], length)
	return s
}

// Text mengembalikan data yang disimpan dalam store sebagai string.
// Fungsi ini mengambil slice byte yang dimulai dari indeks DataStartIndex
// hingga akhir slice dan mengkonversinya menjadi string.
//
// Mengembalikan:
//   - string: Data yang disimpan dalam store, dikonversi dari byte
//     ke string.
func (s Store) Text() string {
	return string(s[DataStartIndex:])
}

// Int mengembalikan data yang disimpan dalam store sebagai int.
// Fungsi ini memeriksa apakah panjang data mencukupi untuk konversi
// ke int. Jika panjang data kurang dari 8 byte, akan mengembalikan
// kesalahan.
//
// Mengembalikan:
//   - int: Data yang disimpan dalam store, dikonversi dari byte
//     ke int.
//   - error: Kesalahan jika panjang data tidak mencukupi untuk
//     konversi.
func (s Store) Int() (int, error) {
	if s.Length() < 8 {
		return 0, fmt.Errorf("insufficient length for int conversion")
	}
	return int(binary.BigEndian.Uint64(s[DataStartIndex:])), nil
}

// Bytes mengembalikan data yang disimpan dalam store sebagai slice byte.
// Fungsi ini mengambil bagian dari store yang dimulai dari indeks
// DataStartIndex hingga akhir, memberikan akses langsung ke data
// mentah yang disimpan.
//
// Mengembalikan:
//   - []byte: Slice byte yang berisi data yang disimpan dalam
//     store, dimulai dari DataStartIndex.
func (s Store) Bytes() []byte {
	return s[DataStartIndex:]
}

// JSON meng-unmarshal data JSON yang disimpan ke dalam struktur tujuan yang diberikan.
// Fungsi ini menggunakan json.Unmarshal untuk mengonversi byte slice
// yang berisi data JSON menjadi tipe data yang ditentukan oleh parameter dest.
//
// Parameter:
//   - dest: Sebuah interface{} yang akan diisi dengan data dari
//     JSON yang disimpan. Struktur tujuan harus cocok dengan format
//     data JSON yang disimpan.
//
// Mengembalikan:
//   - error: Mengembalikan error jika terjadi masalah selama unmarshalling,
//     atau nil jika berhasil.
func (s Store) JSON(dest interface{}) error {
	return json.Unmarshal(s[DataStartIndex:], dest) // Unmarshal data to provided interface{}
}
