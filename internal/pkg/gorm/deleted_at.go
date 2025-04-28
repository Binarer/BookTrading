package gorm

import (
	"database/sql/driver"
	"time"
)

// Scan реализует интерфейс sql.Scanner
func (d *DeletedAt) Scan(value interface{}) error {
	if value == nil {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time, d.Valid = v, true
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		d.Time, d.Valid = t, true
	}
	return nil
}

// Value реализует интерфейс driver.Valuer
func (d DeletedAt) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time, nil
}

// MarshalJSON реализует интерфейс json.Marshaler
func (d DeletedAt) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return d.Time.MarshalJSON()
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (d *DeletedAt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}

	err := d.Time.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	d.Valid = true
	return nil
} 