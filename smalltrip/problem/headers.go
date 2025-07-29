package problem

import (
	"encoding"
	"fmt"
	"net/http"
	"time"
)

type RetryAfter[T time.Time | time.Duration] struct {
	time T
}

var (
	_ encoding.TextMarshaler = RetryAfter[time.Time]{}
	_ fmt.Stringer           = RetryAfter[time.Time]{}
)

func (h RetryAfter[T]) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h RetryAfter[T]) String() string {
	switch t := any(h.time).(type) {
	case time.Time:
		if !t.IsZero() {
			return t.Format(http.TimeFormat)
		}
	case time.Duration:
		return fmt.Sprintf("%.0f", t.Seconds())
	}
	return ""
}
