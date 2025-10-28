package todo

import (
	"time"
)

// Clock interface for time operations, allowing dependency injection for testing.
type Clock interface {
	Now() time.Time
}

// realClock implements Clock using the real time.
type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}
