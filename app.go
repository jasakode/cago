// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago

import (
	"errors"
	"sync"
	"time"
)

// Cago is a lightweight, in‑memory, thread‑safe key/value cache with TTL support.
// It runs a background janitor that periodically removes expired entries.
// All exported methods are safe for concurrent use.
type Cago struct {
	config Config

	// mu protects all fields below.
	mu    sync.RWMutex
	data  map[string]*Entry  // key -> entry
	index map[int64][]string // expiry(unix milli) -> keys

	// lifecycle control for the background janitor
	stop chan struct{}
	done chan struct{}
}

// Config controls cache behavior.
type Config struct {
	// CleanInterval defines how often the janitor scans for expired keys.
	// Default: 1 second if zero.
	CleanInterval time.Duration

	// Timezone is reserved for future use and kept for compatibility.
	// It does not affect cache behavior.
	Timezone Timezone
}

// Global app instance guarded by sync.Once to mimic singleton creation.
var (
	do  sync.Once
	app *Cago
)

// ErrKeyExists is returned by Set when a key already exists and has not expired.
var ErrKeyExists = errors.New("key already exists")

// New initializes the global cache instance and starts the janitor.
// Calling New multiple times is safe; only the first call creates the cache.
func New(conf ...Config) error {
	var cfg Config
	if len(conf) > 0 {
		cfg = conf[0]
	}
	if cfg.CleanInterval <= 0 {
		cfg.CleanInterval = 1 * time.Second
	}
	do.Do(func() {
		app = &Cago{
			config: cfg,
			data:   make(map[string]*Entry),
			index:  make(map[int64][]string),
			stop:   make(chan struct{}),
			done:   make(chan struct{}),
		}
		go app.janitor()
	})
	return nil
}

// Close stops the background janitor and clears the global instance.
func Close() {
	if app == nil {
		return
	}
	app.mu.Lock()
	select {
	case <-app.done:
		// already closed
	default:
		close(app.stop)
	}
	app.mu.Unlock()
	<-app.done
	app = nil
	// allow New to be called again after Close (useful for tests)
	do = sync.Once{}
}

// janitor periodically removes expired entries.
func (c *Cago) janitor() {
	ticker := time.NewTicker(c.config.CleanInterval)
	defer func() {
		ticker.Stop()
		close(c.done)
	}()
	for {
		select {
		case <-ticker.C:
			c.cleanup(time.Now())
		case <-c.stop:
			return
		}
	}
}

// cleanup removes all keys whose expiration time is in the past.
func (c *Cago) cleanup(now time.Time) {
	nowMs := now.UnixMilli()

	c.mu.Lock()
	for exp, keys := range c.index {
		if exp <= nowMs {
			for _, k := range keys {
				if e, ok := c.data[k]; ok && e.isExpiredAt(nowMs) {
					delete(c.data, k)
				}
			}
			delete(c.index, exp)
		}
	}
	c.mu.Unlock()
}

// Set stores a new value for the given key only if the key does not already exist
// or has expired. If the key exists and is not expired, it returns ErrKeyExists.
// ttl <= 0 means the key never expires.
func Set[T any](key string, value T, ttl time.Duration) error {
	if app == nil {
		panic("cago.New must be called before using the cache")
	}

	nowMs := time.Now().UnixMilli()
	var exp int64
	if ttl > 0 {
		exp = nowMs + ttl.Milliseconds()
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	if e, ok := app.data[key]; ok && !e.isExpiredAt(nowMs) {
		return ErrKeyExists
	}

	e := &Entry{
		Key:       key,
		Value:     any(value),
		ExpiresAt: exp,
		CreatedAt: nowMs,
		UpdatedAt: nowMs,
	}
	app.data[key] = e
	if exp > 0 {
		app.index[exp] = append(app.index[exp], key)
	}
	return nil
}

// Put stores or replaces a value for the given key regardless of whether it exists.
// ttl <= 0 means the key never expires.
func Put[T any](key string, value T, ttl time.Duration) {
	if app == nil {
		panic("cago.New must be called before using the cache")
	}

	nowMs := time.Now().UnixMilli()
	var exp int64
	if ttl > 0 {
		exp = nowMs + ttl.Milliseconds()
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	if cur, ok := app.data[key]; ok {
		// remove old index if present
		if cur.ExpiresAt > 0 {
			// no direct reverse index; let janitor naturally remove old exp bucket later
			// because cleaning checks actual expiration time on the entry
		}
	}

	app.data[key] = &Entry{
		Key:       key,
		Value:     any(value),
		ExpiresAt: exp,
		CreatedAt: nowMs,
		UpdatedAt: nowMs,
	}
	if exp > 0 {
		app.index[exp] = append(app.index[exp], key)
	}
}

// Get retrieves a typed value by key. It returns the zero value of T and false
// if the key does not exist or has expired. The type parameter T must match the
// stored value's concrete type.
func Get[T any](key string) (T, bool) {
	var zero T
	if app == nil {
		panic("cago.New must be called before using the cache")
	}

	nowMs := time.Now().UnixMilli()

	app.mu.RLock()
	e, ok := app.data[key]
	app.mu.RUnlock()
	if !ok {
		return zero, false
	}
	if e.isExpiredAt(nowMs) {
		// lazy delete
		Remove(key)
		return zero, false
	}

	v, ok := e.Value.(T)
	if !ok {
		return zero, false
	}
	return v, true
}

// Exist returns true if a non‑expired value exists for key.
func Exist(key string) bool {
	if app == nil {
		panic("cago.New must be called before using the cache")
	}
	nowMs := time.Now().UnixMilli()
	app.mu.RLock()
	e, ok := app.data[key]
	app.mu.RUnlock()
	return ok && !e.isExpiredAt(nowMs)
}

// Remove deletes a key and returns true if it was present.
func Remove(key string) bool {
	if app == nil {
		panic("cago.New must be called before using the cache")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.data[key]; ok {
		delete(app.data, key)
		return true
	}
	return false
}

// Clear deletes all keys from the cache.
func Clear() {
	if app == nil {
		panic("cago.New must be called before using the cache")
	}
	app.mu.Lock()
	app.data = make(map[string]*Entry)
	app.index = make(map[int64][]string)
	app.mu.Unlock()
}
