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
)

// Store is a compact binary blob that prefixes the payload with metadata:
//
//   [0..7]   : createdAt (unix milli)
//   [8..15]  : updatedAt (unix milli)
//   [16..23] : maxAge (user-defined units)
//   [24..31] : payload length (bytes)
//   [32..]   : payload bytes
//
// This format is intended for simple binary serialization of values.
type Store []byte

const (
    CreateAtIndex  = 0
    UpdateAtIndex  = 8
    MaxAgeIndex    = 16
    LengthIndex    = 24
    DataStartIndex = 32
)

// NewStore builds a Store with the given payload and optional maxAge.
// maxAge is a metadata field and not interpreted by Store.
func NewStore(data []byte, maxAge ...uint64) Store {
    MaxAge := uint64(0)
    if len(maxAge) > 0 {
        MaxAge = maxAge[0]
    }

    s := make(Store, DataStartIndex+len(data))
    // createdAt
    binary.BigEndian.PutUint64(s[CreateAtIndex:UpdateAtIndex], uint64(time.Now().UnixMilli()))
    // updatedAt (initially zero)
    for i := UpdateAtIndex; i < MaxAgeIndex; i++ {
        s[i] = 0
    }
    // maxAge
    binary.BigEndian.PutUint64(s[MaxAgeIndex:LengthIndex], MaxAge)
    // payload length
    binary.BigEndian.PutUint64(s[LengthIndex:], uint64(len(data)))
    copy(s[DataStartIndex:], data)
    return s
}

// ParseStore validates the minimum length and returns a Store view over data.
// Returns an empty Store if the input is invalid.
func ParseStore(data []byte) Store {
    if len(data) < DataStartIndex {
        return Store{}
    }
    return Store(data)
}

// Values returns the raw underlying bytes, including metadata and payload.
func (s Store) Values() []byte { return s }

// CreateAt returns the creation timestamp (unix milli).
func (s Store) CreateAt() uint64 {
    return binary.BigEndian.Uint64(s[CreateAtIndex:UpdateAtIndex])
}

// UpdateAt returns the last update timestamp (unix milli).
func (s Store) UpdateAt() uint64 {
    return binary.BigEndian.Uint64(s[UpdateAtIndex:MaxAgeIndex])
}

// SetUpdateAt sets the last update timestamp and returns the mutated store.
func (s Store) SetUpdateAt(date uint64) Store {
    binary.BigEndian.PutUint64(s[UpdateAtIndex:MaxAgeIndex], date)
    return s
}

// Length returns the payload length in bytes. If all=true, returns total store length.
func (s Store) Length(all ...bool) uint64 {
    if len(all) > 0 && all[0] {
        return uint64(len(s))
    }
    return binary.BigEndian.Uint64(s[LengthIndex:])
}

// MaxAge returns the stored maxAge metadata value.
func (s Store) MaxAge() uint64 {
    return binary.BigEndian.Uint64(s[MaxAgeIndex:LengthIndex])
}

// SetMaxAge updates the stored maxAge and returns the mutated store.
func (s Store) SetMaxAge(maxAge uint64) Store {
    binary.BigEndian.PutUint64(s[MaxAgeIndex:LengthIndex], maxAge)
    return s
}

// SetLength updates the payload length and returns the mutated store.
func (s Store) SetLength(length uint64) Store {
    binary.BigEndian.PutUint64(s[LengthIndex:], length)
    return s
}

// Text returns the payload as a string.
func (s Store) Text() string { return string(s[DataStartIndex:]) }

// Int interprets the payload as a big‑endian unsigned 64‑bit integer and
// returns it as int. Returns an error if the payload is too small.
func (s Store) Int() (int, error) {
    if s.Length() < 8 {
        return 0, fmt.Errorf("insufficient length for int conversion")
    }
    return int(binary.BigEndian.Uint64(s[DataStartIndex:])), nil
}

// Bytes returns the payload bytes.
func (s Store) Bytes() []byte { return s[DataStartIndex:] }

// JSON unmarshals the payload into dest.
func (s Store) JSON(dest interface{}) error {
    return json.Unmarshal(s[DataStartIndex:], dest)
}
