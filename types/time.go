package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONTime time.Time
type JSONDate time.Time
type ISODate struct {
	time.Time
	value string
}

// MarshalJSON json time
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}

	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// MarshalJSON json date
func (t JSONDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON ISODate method
func (Date *ISODate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	Date.value = s
	t, _ := time.Parse("2006-01-02", s)
	Date.Time = t
	return nil
}

// MarshalJSON ISODate method
func (Date *ISODate) MarshalJSON() ([]byte, error) {
	return json.Marshal(Date.Time.Format("2006-01-02"))
}

// InRange check if date in range
func (Date *ISODate) InRange(start, end time.Time) bool {
	return Date.Time.After(start) && Date.Time.Before(end)
}

// Value returns value
func (Date *ISODate) Value() string {
	return Date.value
}
