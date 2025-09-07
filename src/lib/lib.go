// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package lib

import (
    "bytes"
    "encoding/binary"
)

// Uint8ToByte converts a uint8 to a 1‑byte slice.
func Uint8ToByte(num uint8) []byte {
    rs := make([]byte, 1)
    rs[0] = num
    return rs
}

// Uint16ToByte converts a uint16 to a 2‑byte big‑endian slice.
func Uint16ToByte(num uint16) []byte {
    rs := make([]byte, 2)
    binary.BigEndian.PutUint16(rs, num)
    return rs
}

// Uint32ToByte converts a uint32 to a 4‑byte big‑endian slice.
func Uint32ToByte(num uint32) []byte {
    rs := make([]byte, 4)
    binary.BigEndian.PutUint32(rs, num)
    return rs
}

// Uint64ToByte converts a uint64 to an 8‑byte big‑endian slice.
func Uint64ToByte(num uint64) []byte {
    rs := make([]byte, 8)
    binary.BigEndian.PutUint64(rs, num)
    return rs
}

// Int8ToByte converts an int8 to a 1‑byte slice using big‑endian encoding.
// For negative values, the two's complement representation is used.
func Int8ToByte(num int8) []byte {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, int8(num)); err != nil {
        panic(err)
    }
    return buf.Bytes()
}

// Int16ToByte converts an int16 to a 2‑byte big‑endian slice.
func Int16ToByte(num int16) []byte {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, int16(num)); err != nil {
        panic(err)
    }
    return buf.Bytes()
}

// Int32ToByte converts an int32 to a 4‑byte big‑endian slice.
func Int32ToByte(num int32) []byte {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, int32(num)); err != nil {
        panic(err)
    }
    return buf.Bytes()
}

// Int64ToByte converts an int64 to an 8‑byte big‑endian slice.
func Int64ToByte(num int64) []byte {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, int64(num)); err != nil {
        panic(err)
    }
    return buf.Bytes()
}

// StringToByte returns the raw bytes of the given string.
func StringToByte(str string) []byte {
    return []byte(str)
}

// StringToByteASCII returns a byte slice of the same length as the input
// string, where non‑ASCII runes (code points > 127) are replaced by 0.
// Note: the length is measured in bytes, not runes, to match Go's string encoding.
func StringToByteASCII(str string) []byte {
    result := make([]byte, len(str))
    for i, c := range str {
        if c > 127 {
            result[i] = 0
        } else {
            result[i] = byte(c)
        }
    }
    return result
}
