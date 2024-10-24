// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"encoding/binary"
)

// Mengubah uint8 ke []byte.
// Fungsi ini akan selalu menghasilkan slice byte dengan panjang 1 byte.
// Nilai yang dapat diubah adalah antara 0 hingga 127.
// Nilai lebih besar dari 127 atau lebih kecil dari 0 akan mengembalikan hasil yang tidak sesuai,
// karena batasan tipe uint8. Misalnya, nilai 128 akan dikonversi menjadi byte dengan nilai 0.
func Uint8ToByte(num uint8) []byte {
	rs := make([]byte, 1)
	rs[0] = num
	return rs
}

// Mengubah uint16 ke []byte.
// Fungsi ini akan selalu menghasilkan slice byte dengan panjang 2 byte.
// Nilai yang dapat diubah adalah antara 0 hingga 65535.
// Nilai di luar batas tipe uint16 akan mengembalikan hasil yang tidak sesuai.
// Misalnya, fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 16-bit ke dalam 2 byte.
func Uint16ToByte(num uint16) []byte {
	rs := make([]byte, 2)
	binary.BigEndian.PutUint16(rs, num)
	return rs
}

// Mengubah uint32 ke []byte.
// Fungsi ini akan selalu menghasilkan slice byte dengan panjang 4 byte.
// Nilai yang dapat diubah adalah antara 0 hingga 4294967295.
// Nilai di luar batas tipe uint32 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 32-bit ke dalam 4 byte.
func Uint32ToByte(num uint32) []byte {
	rs := make([]byte, 4)
	binary.BigEndian.PutUint32(rs, num)
	return rs
}

// Mengubah uint64 ke []byte.
// Fungsi ini akan selalu menghasilkan slice byte dengan panjang 8 byte.
// Nilai yang dapat diubah adalah antara 0 hingga 18446744073709551615.
// Nilai di luar batas tipe uint64 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 64-bit ke dalam 8 byte.
func Uint64ToByte(num uint64) []byte {
	rs := make([]byte, 8)
	binary.BigEndian.PutUint64(rs, num)
	return rs
}

// Mengubah int8 ke []byte.
// Fungsi ini akan menghasilkan slice byte dengan panjang 1 byte.
// Nilai yang dapat diubah adalah antara -128 hingga 127.
// Nilai di luar batas tipe int8 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 8-bit ke dalam 1 byte.
// Jika terjadi kesalahan dalam penulisan buffer, fungsi akan panik.
func Int8ToByte(num int8) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int8(num)) // Konversi int menjadi int64
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Mengubah int16 ke []byte.
// Fungsi ini akan menghasilkan slice byte dengan panjang 2 byte.
// Nilai yang dapat diubah adalah antara -32768 hingga 32767.
// Nilai di luar batas tipe int16 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 16-bit ke dalam 2 byte.
// Jika terjadi kesalahan dalam penulisan buffer, fungsi akan panik.
func Int16ToByte(num int16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int16(num)) // Konversi int menjadi int64
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Mengubah int32 ke []byte.
// Fungsi ini akan menghasilkan slice byte dengan panjang 4 byte.
// Nilai yang dapat diubah adalah antara -2147483648 hingga 2147483647.
// Nilai di luar batas tipe int32 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 32-bit ke dalam 4 byte.
// Jika terjadi kesalahan dalam penulisan buffer, fungsi akan panik.
func Int32ToByte(num int32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(num)) // Konversi int menjadi int64
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Mengubah int64 ke []byte.
// Fungsi ini akan menghasilkan slice byte dengan panjang 8 byte.
// Nilai yang dapat diubah adalah antara -9223372036854775808 hingga 9223372036854775807.
// Nilai di luar batas tipe int64 akan mengembalikan hasil yang tidak sesuai.
// Fungsi ini menggunakan encoding Big Endian untuk menyimpan nilai 64-bit ke dalam 8 byte.
// Jika terjadi kesalahan dalam penulisan buffer, fungsi akan panik.
func Int64ToByte(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int64(num)) // Konversi int menjadi int64
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Mengubah string ke []byte.
// Fungsi ini akan mengembalikan representasi byte dari string yang diberikan
// dengan panjang yang sama dengan string tersebut.
func StringToByte(str string) []byte {
	return []byte(str)
}

// Mengubah string ke []byte dengan batasan ASCII.
// Fungsi ini akan menghasilkan slice byte dengan panjang yang sama dengan string.
// Karakter yang tidak termasuk dalam rentang ASCII (0-127) akan diubah menjadi null (0).
// Ini memastikan bahwa hasilnya hanya berisi karakter-karakter ASCII.
func StringToByteASCII(str string) []byte {
	// Buat slice byte dengan panjang sama dengan string
	result := make([]byte, len(str))
	for i, c := range str {
		// Pastikan karakter adalah ASCII
		if c > 127 {
			// Jika bukan, masukkan karakter null atau bisa juga ditangani dengan cara lain
			result[i] = 0
		} else {
			result[i] = byte(c)
		}
	}
	return result
}
