package cago

// Entry represents a single cache item.
// Value holds the stored data, ExpiresAt is a unix timestamp in milliseconds.
// A zero ExpiresAt means the entry never expires.
type Entry struct {
	Key       string
	Value     any
	ExpiresAt int64 // unix milli; 0 means never expires
	CreatedAt int64 // unix milli
	UpdatedAt int64 // unix milli
}

// isExpiredAt reports whether the entry should be considered expired at the
// given unix milli timestamp.
func (e *Entry) isExpiredAt(nowMs int64) bool {
	return e.ExpiresAt > 0 && nowMs >= e.ExpiresAt
}
