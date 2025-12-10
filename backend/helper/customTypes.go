package helper

import (
	"strings"
	"time"
)

// Date is a custom type that accepts both RFC3339 and YYYY-MM-DD formats
type Date struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}

	// Try RFC3339 format first
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try date-only format (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		d.Time = t
		return nil
	}

	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
// This is used by form binding (multipart/form-data and application/x-www-form-urlencoded)
func (d *Date) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	// Try RFC3339 format first
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try date-only format (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		d.Time = t
		return nil
	}

	return err
}

// MarshalJSON implements the json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + d.Time.Format(time.RFC3339) + "\""), nil
}
