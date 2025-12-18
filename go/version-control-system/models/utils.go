package models

import "time"

// GetCurrentTime returns the current time (useful for testing)
var GetCurrentTime = func() time.Time {
	return time.Now()
}

