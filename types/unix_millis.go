package types

import "time"

// Milliseconds represents a time or duration in milliseconds.
type Milliseconds int64

// Time returns the time associated with the Unix milliseconds.
func (m Milliseconds) Time() time.Time {
	return time.Unix(0, int64(m)*int64(time.Millisecond/time.Nanosecond))
}

// Duration returns the duration with the given number of milliseconds.
func (m Milliseconds) Duration() time.Duration {
	return time.Millisecond * time.Duration(m)
}
